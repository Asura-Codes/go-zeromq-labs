package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab20/internal/tracing"
	"github.com/go-zeromq/zmq4"
)

func main() {
	port := flag.Int("port", 5555, "Port to listen for spans")
	flag.Parse()

	log.Printf("Trace Collector starting on :%d...", *port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down collector...")
		cancel()
	}()

	pull := zmq4.NewPull(ctx)
	defer pull.Close()

	addr := fmt.Sprintf("tcp://*:%d", *port)
	if err := pull.Listen(addr); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Collector ready to receive spans.")

	for {
		msg, err := pull.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("Recv error: %v", err)
			continue
		}

		var span tracing.Span
		if err := json.Unmarshal(msg.Frames[0], &span); err != nil {
			log.Printf("Failed to unmarshal span: %v", err)
			continue
		}

		log.Printf("[TRACE:%s] %s -> %s (%s) [%s]", 
			span.TraceID, span.ServiceName, span.Operation, span.Duration, span.Status)
	}
}
