package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-zeromq/zmq4"
)

func main() {
	brokerAddr := flag.String("broker", "tcp://localhost:5556", "Address of the broker backend")
	identity := flag.String("id", "", "Identity of the worker")
	flag.Parse()

	if *identity == "" {
		hostname, _ := os.Hostname()
		*identity = fmt.Sprintf("worker-%s-%d", hostname, os.Getpid())
	}

	ctx := context.Background()
	worker := zmq4.NewRep(ctx, zmq4.WithID(zmq4.SocketIdentity(*identity)))
	defer worker.Close()

	log.Printf("Worker [%s] connecting to %s", *identity, *brokerAddr)
	if err := worker.Dial(*brokerAddr); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	for {
		msg, err := worker.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
			break
		}

		// Simulate some "work"
		payload := string(msg.Frames[0])
		// log.Printf("[%s] Received: %s", *identity, payload)

		reply := fmt.Sprintf("ACK from %s to: %s", *identity, payload)
		if err := worker.Send(zmq4.NewMsgFrom([]byte(reply))); err != nil {
			log.Printf("Send error: %v", err)
			break
		}
		
		// To keep the console from being TOO messy but still show action
		if time.Now().Unix() % 2 == 0 {
			fmt.Print(".")
		}
	}
}
