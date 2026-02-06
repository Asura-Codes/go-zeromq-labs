package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"gemini-zeromq-labs/lab14/internal/dht"

	"github.com/go-zeromq/zmq4"
)

type Node struct {
	Addr  string
	Ring  *dht.VirtualRing
	Store map[string]string
	mu    sync.RWMutex
}

func main() {
	port := flag.Int("port", 5001, "Port to listen on")
	peers := flag.String("peers", "", "Comma-separated list of all nodes in the ring (including self)")
	vnodeCount := flag.Int("vnodes", 10, "Number of virtual nodes per physical node")
	flag.Parse()

	addr := fmt.Sprintf("tcp://127.0.0.1:%d", *port)
	peerList := strings.Split(*peers, ",")

	ring := &dht.VirtualRing{}
	for _, p := range peerList {
		if p != "" {
			ring.AddNode(p, *vnodeCount)
		}
	}

	n := &Node{
		Addr:  addr,
		Ring:  ring,
		Store: make(map[string]string),
	}

	log.Printf("Starting DHT Node %s with %d vnodes. Total ring points: %d", n.Addr, *vnodeCount, len(ring.VNodes))

	ctx := context.Background()
	router := zmq4.NewRouter(ctx)
	defer router.Close()

	if err := router.Listen(n.Addr); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	for {
		msg, err := router.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
			continue
		}

		frames := msg.Frames
		if len(frames) < 3 {
			continue
		}

		clientID := frames[0]
		// frames[1] empty
		command := string(frames[2])

		switch command {
		case "PUT":
			if len(frames) < 5 {
				continue
			}
			key := string(frames[3])
			val := string(frames[4])

			target := n.Ring.GetResponsibleNode(key)
			if target == n.Addr {
				n.mu.Lock()
				n.Store[key] = val
				n.mu.Unlock()
				log.Printf("[LOCAL] Stored: %s = %s", key, val)
				router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, []byte("OK")}})
			} else {
				log.Printf("[PROXY] Forwarding PUT %s to node %s", key, target)
				res, err := n.proxyRequest(ctx, target, frames[2:])
				if err == nil {
					router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, res.Frames[0]}})
				} else {
					router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, []byte("ERROR")}})
				}
			}

		case "GET":
			if len(frames) < 4 {
				continue
			}
			key := string(frames[3])

			target := n.Ring.GetResponsibleNode(key)
			if target == n.Addr {
				n.mu.RLock()
				val, ok := n.Store[key]
				n.mu.RUnlock()

				if ok {
					log.Printf("[LOCAL] Retrieval: %s", key)
					router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, []byte("OK"), []byte(val)}})
				} else {
					router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, []byte("NOT_FOUND")}})
				}
			} else {
				log.Printf("[PROXY] Forwarding GET %s to node %s", key, target)
				res, err := n.proxyRequest(ctx, target, frames[2:])
				if err == nil {
					// Prepend client routing frames
					reply := append([][]byte{clientID, {}}, res.Frames...)
					router.Send(zmq4.Msg{Frames: reply})
				} else {
					router.Send(zmq4.Msg{Frames: [][]byte{clientID, {}, []byte("NOT_FOUND")}})
				}
			}
		}
	}
}

func (n *Node) proxyRequest(ctx context.Context, target string, payload [][]byte) (zmq4.Msg, error) {
	client := zmq4.NewReq(ctx)
	defer client.Close()
	if err := client.Dial(target); err != nil {
		return zmq4.Msg{}, err
	}

	err := client.Send(zmq4.Msg{Frames: payload})
	if err != nil {
		return zmq4.Msg{}, err
	}

	return client.Recv()
}
