package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Initialize ZeroMQ PUB socket
	// Agent Connects, Dashboard Binds.
	pub := zmq4.NewPub(ctx)
	defer pub.Close()

	logger.Info("Agent connecting", "endpoint", cfg.Endpoint)
	if err := pub.Dial(cfg.Endpoint); err != nil {
		logger.Error("Failed to connect to endpoint", "error", err)
		os.Exit(1)
	}

	hostname, _ := os.Hostname()
	serviceName := "monitor-agent-01"

	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer ticker.Stop()

	logger.Info("Starting metrics publishing", "topic", cfg.Topic, "interval_sec", cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down agent...")
			return
		case <-ticker.C:
			// Generate mock metrics
			metric := protocol.Metric{
				Timestamp: time.Now(),
				Service:   serviceName,
				Host:      hostname,
				CPU:       rand.Float64() * 100,
				Memory:    rand.Float64() * 16384, // MB
				Status:    "OK",
			}

			payload, err := metric.ToJSON()
			if err != nil {
				logger.Error("Error serializing metric", "error", err)
				continue
			}

			// Create ZMQ message: [Topic] [Payload]
			msg := zmq4.NewMsgFrom(
				[]byte(cfg.Topic),
				payload,
			)

			if err := pub.Send(msg); err != nil {
				logger.Error("Error sending message", "error", err)
			} else {
				logger.Info("Sent metric", "cpu", metric.CPU, "memory_mb", metric.Memory)
			}
		}
	}
}
