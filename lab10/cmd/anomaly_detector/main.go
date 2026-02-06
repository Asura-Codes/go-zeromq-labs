package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"gemini-zeromq-labs/lab10/internal/config"
	"gemini-zeromq-labs/lab10/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting Anomaly Detector...")

	rep := zmq4.NewRep(context.Background())
	defer rep.Close()

	if err := rep.Listen(config.AnomalyRepAddress); err != nil {
		log.Fatalf("Failed to bind REP socket: %v", err)
	}
	log.Printf("Listening for scan requests on %s", config.AnomalyRepAddress)

	for {
		msg, err := rep.Recv()
		if err != nil {
			log.Printf("Error receiving: %v", err)
			continue
		}

		ip := string(msg.Bytes())
		log.Printf("Scanning IP: %s", ip)

		// Simulate processing time
		time.Sleep(200 * time.Millisecond)

		// Random verdict
		status := protocol.StatusSafe
		if rand.Float32() < 0.3 { // 30% chance of malicious
			status = protocol.StatusMalicious
		}

		if err := rep.Send(zmq4.NewMsgString(status)); err != nil {
			log.Printf("Error sending reply: %v", err)
		}
	}
}
