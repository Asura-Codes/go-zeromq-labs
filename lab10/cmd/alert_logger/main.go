package main

import (
	"context"
	"log"

	"gemini-zeromq-labs/lab10/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting Alert Logger...")

	pull := zmq4.NewPull(context.Background())
	defer pull.Close()

	if err := pull.Listen(config.AlertPullAddress); err != nil {
		log.Fatalf("Failed to bind PULL socket: %v", err)
	}
	log.Printf("Listening for alerts on %s", config.AlertPullAddress)

	for {
		msg, err := pull.Recv()
		if err != nil {
			log.Printf("Error receiving alert: %v", err)
			continue
		}

		log.Printf("[CRITICAL ALERT] %s", string(msg.Bytes()))
	}
}
