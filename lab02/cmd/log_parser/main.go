package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab02/internal/config"
	"gemini-zeromq-labs/lab02/internal/protocol"

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

	// 1. Socket to receive messages on (PULL) from Collector
	receiver := zmq4.NewPull(ctx)
	defer receiver.Close()
	
	collectorAddr := cfg.CollectorConnectAddr()
	logger.Info("Worker connecting to collector", "endpoint", collectorAddr)
	if err := receiver.Dial(collectorAddr); err != nil {
		logger.Error("Failed to dial collector", "error", err)
		os.Exit(1)
	}

	// 2. Socket to send messages to (PUSH) to Sink
	sender := zmq4.NewPush(ctx)
	defer sender.Close()

	sinkAddr := cfg.SinkConnectAddr()
	logger.Info("Worker connecting to sink", "endpoint", sinkAddr)
	if err := sender.Dial(sinkAddr); err != nil {
		logger.Error("Failed to dial sink", "error", err)
		os.Exit(1)
	}

	// Regex for basic IP redaction
	ipRegex := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	logger.Info("Worker started. Waiting for tasks...")

	for {
		// Read raw message
		msg, err := receiver.Recv()
		if err != nil {
			// Context canceled or error
			if ctx.Err() != nil {
				break
			}
			logger.Error("Worker error receiving", "error", err)
			continue
		}

		// Unmarshal
		var entry protocol.LogEntry
		if err := json.Unmarshal(msg.Bytes(), &entry); err != nil {
			logger.Error("Worker error unmarshalling", "error", err)
			continue
		}

		// Process / Anonymize
		cleanMsg := ipRegex.ReplaceAllString(entry.Message, "[REDACTED]")
		sanitized := cleanMsg != entry.Message

		processed := protocol.ProcessedLogEntry{
			OriginalID:   entry.ID,
			ProcessedAt:  time.Now(),
			Level:        entry.Level,
			Sanitized:    sanitized,
			CleanMessage: cleanMsg,
		}

		outData, err := json.Marshal(processed)
		if err != nil {
			logger.Error("Worker error marshalling processed data", "error", err)
			continue
		}

		// Send to Sink
		if err := sender.Send(zmq4.NewMsg(outData)); err != nil {
			logger.Error("Worker error sending to sink", "error", err)
		}
	}
}
