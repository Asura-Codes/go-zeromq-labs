package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"gemini-zeromq-labs/lab18/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Camera Feed Starting...")

	pub := zmq4.NewPub(context.Background())
	if err := pub.Listen(config.StreamAddr); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	frameID := uint64(0)
	// Pre-allocate frame buffer to minimize GC pressure (simulating zero-copy intent)
	buffer := make([]byte, config.FrameSize)

	ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
	defer ticker.Stop()

	for range ticker.C {
		frameID++
		
		// Simulate data change
		buffer[0] = byte(frameID % 256)
		buffer[1] = byte(rand.Intn(256))

		// In go-zeromq, Msg is a slice of frames.
		// We send the ID as frame 0 and the raw data as frame 1.
		msg := zmq4.Msg{
			Frames: [][]byte{
				[]byte{byte(frameID >> 24), byte(frameID >> 16), byte(frameID >> 8), byte(frameID)},
				buffer,
			},
		}

		if err := pub.Send(msg); err != nil {
			log.Printf("Send error: %v", err)
		}

		if frameID%30 == 0 {
			log.Printf("Published 30 frames. Latest ID: %d", frameID)
		}
	}
}
