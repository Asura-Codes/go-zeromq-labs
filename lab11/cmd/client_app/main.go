package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gemini-zeromq-labs/lab11/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting Client Application...")
	
	endpoints := []string{config.PrimaryFrontendAddr, config.BackupFrontendAddr}
	
	// Create context
	ctx := context.Background()
	req := zmq4.NewReq(ctx)
	defer req.Close()

	// Connect to both potentially, or strategy: connect to one, if fail, connect to other.
	// For simplicity, let's just connect to both and ZMQ will load balance or we manage it.
	// A better HA client strategy:
	// 1. Try Primary. 
	// 2. If timeout, Try Backup.
	
	// Let's implement an explicit failover loop
	activeEndpoint := 0
	
	// Initial dial
	log.Printf("Connecting to %s...", endpoints[activeEndpoint])
	if err := req.Dial(endpoints[activeEndpoint]); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	seq := 0

	for range ticker.C {
		seq++
		msgContent := fmt.Sprintf("Hello %d", seq)
		log.Printf("Sending: %s", msgContent)

		// Send
		if err := req.Send(zmq4.NewMsgString(msgContent)); err != nil {
			log.Printf("Send failed: %v", err)
			reconnect(&req, endpoints, &activeEndpoint)
			continue
		}

		// Wait for reply with timeout using goroutine/select pattern
		type recvResult struct {
			msg zmq4.Msg
			err error
		}
		resultChan := make(chan recvResult, 1)

		go func() {
			msg, err := req.Recv()
			resultChan <- recvResult{msg, err}
		}()

		select {
		case res := <-resultChan:
			if res.err != nil {
				log.Printf("Recv error: %v", res.err)
				reconnect(&req, endpoints, &activeEndpoint)
			} else {
				log.Printf("Reply: %s", string(res.msg.Bytes()))
			}
		case <-time.After(2 * time.Second):
			log.Println("Timeout! Switching server...")
			// Closing the socket will cause the blocked Recv to error out/return
			reconnect(&req, endpoints, &activeEndpoint)
		}
	}
}

func reconnect(sock *zmq4.Socket, endpoints []string, index *int) {
	// Close current (actually ZMQ REQ sockets are tricky to reuse if request in flight, 
	// usually better to destroy and recreate socket on hard failover)
	(*sock).Close()
	
	newSock := zmq4.NewReq(context.Background())
	*sock = newSock

	// Toggle index
	*index = (*index + 1) % len(endpoints)
	endpoint := endpoints[*index]
	
	log.Printf("Reconnecting to %s...", endpoint)
	if err := (*sock).Dial(endpoint); err != nil {
		log.Printf("Failed to dial: %v", err)
	}
}
