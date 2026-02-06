package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab01/internal/config"
	"gemini-zeromq-labs/lab01/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create context that cancels on interrupt signals

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Initialize ZeroMQ SUB socket
	sub := zmq4.NewSub(ctx)
	defer sub.Close()

	logger.Info("Dashboard listening", "endpoint", cfg.Endpoint)
	if err := sub.Listen(cfg.Endpoint); err != nil {
		logger.Error("Failed to bind to endpoint", "error", err)
		os.Exit(1)
	}

	// Subscribe to topic
	logger.Info("Subscribing to topic", "topic", cfg.Topic)
	if err := sub.SetOption(zmq4.OptionSubscribe, cfg.Topic); err != nil {
		logger.Error("Failed to subscribe", "error", err)
		os.Exit(1)
	}

	logger.Info("Waiting for metrics...")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down dashboard...")
			return
		default:
			// Read message
			msg, err := sub.Recv()
			if err != nil {
				// check if context cancelled
				if ctx.Err() != nil {
					return
				}
				logger.Error("Error receiving message", "error", err)
				continue
			}

			// Expected Frames: [Topic, Payload]
			if len(msg.Frames) < 2 {
				logger.Warn("Received malformed message", "frames", len(msg.Frames))
				continue
			}

			// topic := string(msg.Frames[0])
			payload := msg.Frames[1]

			metric, err := protocol.FromJSON(payload)
			if err != nil {
				logger.Error("Error parsing metric", "error", err)
				continue
			}

			// Print formatted (keep fmt for Dashboard display as it is a UI)
			fmt.Printf("[%s] %s@%s | CPU: %.2f%% | Mem: %.0f MB | Status: %s\n",
				metric.Timestamp.Format("15:04:05"),
				metric.Service,
				metric.Host,
				metric.CPU,
				metric.Memory,
				metric.Status,
			)
		}
	}
}
