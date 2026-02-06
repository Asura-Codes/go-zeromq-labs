package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab04/internal/config"
	"gemini-zeromq-labs/lab04/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	topic := "sensors/temp"
	if len(os.Args) > 1 {
		topic = os.Args[1]
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	// Connect to Broker Frontend
	sub := zmq4.NewSub(ctx)
	defer sub.Close()

	addr := cfg.SubConnectAddr()
	logger.Info("Terminal connecting to broker", "endpoint", addr)
	if err := sub.Dial(addr); err != nil {
		logger.Error("Failed to connect", "error", err)
		os.Exit(1)
	}

	logger.Info("Subscribing", "topic", topic)
	if err := sub.SetOption(zmq4.OptionSubscribe, topic); err != nil {
		logger.Error("Failed to subscribe", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Listening for updates on %s...\n", topic)

	for {
		msg, err := sub.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Error("Recv error", "error", err)
			continue
		}

		if len(msg.Frames) < 2 {
			continue
		}

		// rcvTopic := string(msg.Frames[0])
		payload := msg.Frames[1]

		data, err := protocol.FromBytes(payload)
		if err != nil {
			logger.Error("Parse error", "error", err)
			continue
		}

		fmt.Printf("[%s] %s: %.2f (cached or new)\n", 
			data.Timestamp.Format("15:04:05.000"), 
			data.SensorID,
			data.Value,
		)
	}
}
