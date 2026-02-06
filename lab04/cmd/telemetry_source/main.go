package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab04/internal/config"
	"gemini-zeromq-labs/lab04/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Connect to Broker Backend
	pub := zmq4.NewPub(ctx)
	defer pub.Close()

	addr := cfg.PubConnectAddr()
	logger.Info("Source connecting to broker", "endpoint", addr)
	if err := pub.Dial(addr); err != nil {
		logger.Error("Failed to connect", "error", err)
		os.Exit(1)
	}

	// Sensors
	sensors := []string{"sensors/temp", "sensors/pressure", "sensors/humidity"}
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logger.Info("Publishing telemetry...")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Pick random sensor
			topic := sensors[rand.Intn(len(sensors))]
			
			// Generate data
			data := protocol.TelemetryData{
				SensorID:  topic,
				Value:     rand.Float64() * 100,
				Timestamp: time.Now(),
				Unit:      "raw",
			}
			
			payload, _ := data.ToBytes()

			// Send [Topic] [Payload]
			msg := zmq4.NewMsgFrom(
				[]byte(topic),
				payload,
			)

			if err := pub.Send(msg); err != nil {
				logger.Error("Send error", "error", err)
			} else {
				logger.Info("Published", "topic", topic, "val", data.Value)
			}
		}
	}
}
