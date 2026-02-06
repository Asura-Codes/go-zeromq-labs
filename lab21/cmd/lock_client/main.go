package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"gemini-zeromq-labs/lab21/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

func main() {
	id := flag.String("id", fmt.Sprintf("client-%d", rand.Intn(1000)), "Client ID")
	resource := flag.String("resource", "shared-db", "Resource to lock")
	serversStr := flag.String("servers", "tcp://127.0.0.1:5555,tcp://127.0.0.1:5556,tcp://127.0.0.1:5557", "Comma-separated server addresses")
	flag.Parse()

	serverAddrs := strings.Split(*serversStr, ",")
	log.Printf("[%s] Starting with %d servers...", *id, len(serverAddrs))

	ctx := context.Background()
	var sockets []zmq4.Socket
	for _, addr := range serverAddrs {
		dealer := zmq4.NewDealer(ctx)
		if err := dealer.Dial(addr); err != nil {
			log.Printf("Warning: Failed to dial %s: %v", addr, err)
			continue
		}
		sockets = append(sockets, dealer)
	}
	defer func() {
		for _, s := range sockets {
			s.Close()
		}
	}()

	quorum := (len(sockets) / 2) + 1

	for {
		log.Printf("[%s] Attempting to acquire quorum lock on %s (Quorum size: %d)...", *id, *resource, quorum)
		
		grantedIndices := acquireQuorumLock(sockets, *id, *resource, quorum)
		if len(grantedIndices) >= quorum {
			log.Printf("[%s] QUORUM REACHED (%d/%d). Doing work...", *id, len(grantedIndices), len(sockets))
			
			done := make(chan bool)
			go heartbeatQuorum(sockets, grantedIndices, *id, *resource, done)
			
			workTime := time.Duration(3+rand.Intn(5)) * time.Second
			time.Sleep(workTime)
			
			done <- true
			releaseQuorumLock(sockets, grantedIndices, *id, *resource)
			log.Printf("[%s] Work complete, quorum lock released.", *id)
			
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("[%s] Quorum failed (%d/%d). Retrying...", *id, len(grantedIndices), len(sockets))
			// Release any partial locks acquired
			releaseQuorumLock(sockets, grantedIndices, *id, *resource)
			time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
		}
	}
}

func acquireQuorumLock(sockets []zmq4.Socket, clientID, resource string, quorum int) []int {
	var granted []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	req := protocol.LockRequest{
		Type:     protocol.Acquire,
		Resource: resource,
		ClientID: clientID,
		TTL:      10, // Longer TTL for quorum
	}
	data, _ := json.Marshal(req)

	for i, s := range sockets {
		wg.Add(1)
		go func(idx int, sock zmq4.Socket) {
			defer wg.Done()
			sock.Send(zmq4.NewMsgFrom([]byte(""), data))
			
			// Set a short timeout for the response
			msg, err := sock.Recv()
			if err == nil {
				var resp protocol.LockResponse
				json.Unmarshal(msg.Frames[1], &resp)
				if resp.Type == protocol.Grant {
					mu.Lock()
					granted = append(granted, idx)
					mu.Unlock()
				}
			}
		}(i, s)
	}

	wg.Wait()
	return granted
}

func heartbeatQuorum(sockets []zmq4.Socket, indices []int, clientID, resource string, done chan bool) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			req := protocol.LockRequest{
				Type:     protocol.Heartbeat,
				Resource: resource,
				ClientID: clientID,
				TTL:      10,
			}
			data, _ := json.Marshal(req)
			for _, idx := range indices {
				sockets[idx].Send(zmq4.NewMsgFrom([]byte(""), data))
				sockets[idx].Recv() // Clear response
			}
		}
	}
}

func releaseQuorumLock(sockets []zmq4.Socket, indices []int, clientID, resource string) {
	req := protocol.LockRequest{
		Type:     protocol.Release,
		Resource: resource,
		ClientID: clientID,
	}
	data, _ := json.Marshal(req)
	for _, idx := range indices {
		sockets[idx].Send(zmq4.NewMsgFrom([]byte(""), data))
		sockets[idx].Recv()
	}
}