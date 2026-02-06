package main

import (
	"context"
	"log"
	"time"

	"gemini-zeromq-labs/lab18/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Analytics Engine Starting...")

	sub := zmq4.NewSub(context.Background())
	if err := sub.Dial(config.StreamAddr); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	sub.SetOption(zmq4.OptionSubscribe, "")

	var totalBytes uint64
	frameCount := 0
	start := time.Now()

	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
			continue
		}

		if len(msg.Frames) < 2 {
			continue
		}

		// Simulate "Processing"
		frameCount++
		totalBytes += uint64(len(msg.Frames[1]))

		if frameCount%30 == 0 {
			elapsed := time.Since(start).Seconds()
			throughput := float64(totalBytes) / (1024 * 1024 * elapsed)
			log.Printf("Processed 30 frames. Total Bytes: %d MB. Throughput: %.2f MB/s", totalBytes/(1024*1024), throughput)
		}
	}
}
