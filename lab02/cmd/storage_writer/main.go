package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab02/internal/config"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		logger.Info("Interrupt received, shutting down...")
		cancel()
	}()

	// 1. Create PULL socket
	sink := zmq4.NewPull(ctx)
	defer sink.Close()

	// 2. Bind to endpoint
	bindAddr := cfg.SinkBindAddr()
	logger.Info("Sink binding", "endpoint", bindAddr)
	if err := sink.Listen(bindAddr); err != nil {
		logger.Error("Failed to bind sink", "error", err)
		os.Exit(1)
	}

	logger.Info("Sink ready. Waiting for processed logs...")

	var count int
	start := time.Now()
	first := true

	for {
		_, err := sink.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Error("Sink error receiving", "error", err)
			continue
		}

		if first {
			start = time.Now()
			first = false
			logger.Info("Batch started")
		}

		count++

		if count%1000 == 0 {
			// Using fmt here just for progress visibility in console if needed, or logger
			fmt.Printf("\rSink: Processed %d messages...", count)
		}
	}

	duration := time.Since(start)
	logger.Info("Sink finished", "total", count, "duration", duration.String())
}
