package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-zeromq/zmq4"
)

func main() {
	frontendAddr := flag.String("frontend", "tcp://*:5555", "Address to bind for clients (ROUTER)")
	backendAddr := flag.String("backend", "tcp://*:5556", "Address to bind for workers (DEALER)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	log.Printf("Starting High-Performance Central Message Hub...")

	frontend := zmq4.NewRouter(ctx)
	defer frontend.Close()
	if err := frontend.Listen(*frontendAddr); err != nil {
		log.Fatalf("Failed to bind frontend: %v", err)
	}

	backend := zmq4.NewDealer(ctx)
	defer backend.Close()
	if err := backend.Listen(*backendAddr); err != nil {
		log.Fatalf("Failed to bind backend: %v", err)
	}

	var msgCount uint64

	// Stats reporter
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		var lastCount uint64
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				currentCount := atomic.LoadUint64(&msgCount)
				diff := currentCount - lastCount
				log.Printf("[STATS] Throughput: %.2f msg/s (Total: %d)", float64(diff)/2.0, currentCount)
				lastCount = currentCount
			}
		}
	}()

	// Forward: Frontend -> Backend
	go func() {
		for {
			msg, err := frontend.Recv()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			atomic.AddUint64(&msgCount, 1)
			if err := backend.Send(msg); err != nil {
				log.Printf("[ERROR] Backend Send: %v", err)
			}
		}
	}()

	// Forward: Backend -> Frontend
	go func() {
		for {
			msg, err := backend.Recv()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			if err := frontend.Send(msg); err != nil {
				log.Printf("[ERROR] Frontend Send: %v", err)
			}
		}
	}()

	log.Printf("Hub is active. Routing traffic between %s and %s", *frontendAddr, *backendAddr)
	<-ctx.Done()
	log.Printf("Shutting down...")
}