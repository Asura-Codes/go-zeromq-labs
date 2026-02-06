package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"gemini-zeromq-labs/lab16/internal/config"
	"github.com/go-zeromq/zmq4"
)

type Event struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
}

func main() {
	log.Printf("Starting Telemetry Producer at %s", config.ZmqPubAddr)

	pub := zmq4.NewPub(context.Background())
	defer pub.Close()

	if err := pub.Listen(config.ZmqPubAddr); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	eventTypes := []string{"LOGIN", "FILE_ACCESS", "NETWORK_CON", "SUDO"}
	severities := []string{"INFO", "WARNING", "CRITICAL"}

	for {
		event := Event{
			Timestamp: time.Now().Format(time.RFC3339),
			Type:      eventTypes[rand.Intn(len(eventTypes))],
			Severity:  severities[rand.Intn(len(severities))],
			Message:   "Simulated security event",
		}

		data, _ := json.Marshal(event)
		log.Printf("Broadcasting: %s", string(data))

		err := pub.Send(zmq4.Msg{Frames: [][]byte{[]byte("events"), data}})
		if err != nil {
			log.Printf("Send error: %v", err)
		}

		time.Sleep(2 * time.Second)
	}
}
