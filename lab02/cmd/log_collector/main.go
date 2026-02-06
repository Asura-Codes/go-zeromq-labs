package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab02/internal/config"
	"gemini-zeromq-labs/lab02/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// 1. Create a PUSH socket (Ventilator)
	ventilator := zmq4.NewPush(ctx)
	defer ventilator.Close()

	// 2. Bind to the endpoint
	bindAddr := cfg.CollectorBindAddr()
	logger.Info("Collector binding", "endpoint", bindAddr)
	if err := ventilator.Listen(bindAddr); err != nil {
		logger.Error("Failed to bind ventilator", "error", err)
		os.Exit(1)
	}

	// 3. Wait for workers to connect
	logger.Info("Collector ready. Waiting 5 seconds for workers to connect...")
	time.Sleep(5 * time.Second)

	logger.Info("Starting log generation...")
	start := time.Now()

	// Generate a batch of logs
	const batchSize = 10000
	sources := []string{"Firewall-1", "Auth-Server", "Web-Gateway", "DB-Cluster"}
	levels := []protocol.LogLevel{protocol.INFO, protocol.WARN, protocol.ERROR, protocol.DEBUG}

	for i := 0; i < batchSize; i++ {
		// Check for cancellation
		if ctx.Err() != nil {
			break
		}

		// Simulate some "IP" in the message to be anonymized later
		msgContent := fmt.Sprintf("User login attempt from IP: 192.168.1.%d", rand.Intn(255))
		if rand.Float32() < 0.1 {
			msgContent = fmt.Sprintf("Malware detected from IP: 10.0.0.%d", rand.Intn(255))
		}

		entry := protocol.LogEntry{
			ID:        fmt.Sprintf("log-%d", i),
			Timestamp: time.Now(),
			Level:     levels[rand.Intn(len(levels))],
			Source:    sources[rand.Intn(len(sources))],
			Message:   msgContent,
			RawData:   "0xDEADBEEF...", // filler
		}

		data, err := json.Marshal(entry)
		if err != nil {
			logger.Error("Error marshalling log", "error", err)
			continue
		}

		// Send task
		msg := zmq4.NewMsg(data)
		if err := ventilator.Send(msg); err != nil {
			logger.Error("Failed to send log", "error", err)
		}
	}

	elapsed := time.Since(start)
	logger.Info("Batch completed", "count", batchSize, "duration", elapsed.String())

	// Give time for packets to flush out before closing
	time.Sleep(1 * time.Second)
}
