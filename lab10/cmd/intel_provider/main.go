package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gemini-zeromq-labs/lab10/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting Intel Provider...")

	// Create a PUB socket
	pub := zmq4.NewPub(context.Background())
	defer pub.Close()

	if err := pub.Listen(config.IntelPubAddress); err != nil {
		log.Fatalf("Failed to bind PUB socket: %v", err)
	}
	log.Printf("Publishing Threat Intel on %s", config.IntelPubAddress)

	// Simulate generating IPs
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ip := fmt.Sprintf("192.168.1.%d", rand.Intn(255))
		log.Printf("Broadcasting suspect IP: %s", ip)
		
		msg := zmq4.NewMsgString(ip)
		if err := pub.Send(msg); err != nil {
			log.Printf("Failed to send: %v", err)
		}
	}
}
