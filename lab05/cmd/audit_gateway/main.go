package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab05/internal/config"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Audit Gateway (ROUTER-DEALER)...")

	// Context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Signal received, shutting down...")
		cancel()
	}()

	// 1. Prepare Frontend (ROUTER) for Clients
	frontend := zmq4.NewRouter(ctx)
	defer frontend.Close()
	if err := frontend.Listen(config.GatewayFrontendAddr); err != nil {
		logger.Error("Failed to listen on frontend", "addr", config.GatewayFrontendAddr, "error", err)
		os.Exit(1)
	}
	logger.Info("Frontend (ROUTER) listening", "addr", config.GatewayFrontendAddr)

	// 2. Prepare Backend (DEALER) for Workers
	backend := zmq4.NewDealer(ctx)
	defer backend.Close()
	if err := backend.Listen(config.GatewayBackendAddr); err != nil {
		logger.Error("Failed to listen on backend", "addr", config.GatewayBackendAddr, "error", err)
		os.Exit(1)
	}
	logger.Info("Backend (DEALER) listening", "addr", config.GatewayBackendAddr)

	// 3. Start Proxy
	logger.Info("Gateway active. Proxying messages...")

	// Channels to bridge messages
	// We use a custom struct to identify source
	type bridgeMsg struct {
		msg  zmq4.Msg
		from string // "frontend" or "backend"
	}

	msgChan := make(chan bridgeMsg)

	// Frontend Reader
	go func() {
		for {
			msg, err := frontend.Recv()
			if err != nil {
				// Often "resource temporarily unavailable" or context cancel
				return
			}
			select {
			case msgChan <- bridgeMsg{msg, "frontend"}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Backend Reader
	go func() {
		for {
			msg, err := backend.Recv()
			if err != nil {
				return
			}
			select {
			case msgChan <- bridgeMsg{msg, "backend"}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Main Loop: Serializes writes
	for {
		select {
		case m := <-msgChan:
			if m.from == "frontend" {
				// Route to Backend
				if err := backend.Send(m.msg); err != nil {
					logger.Error("Failed to forward to backend", "error", err)
				}
			} else {
				// Route to Frontend
				if err := frontend.Send(m.msg); err != nil {
					logger.Error("Failed to forward to frontend", "error", err)
				}
			}
		case <-ctx.Done():
			goto Exit
		}
	}

Exit:
	logger.Info("Gateway stopped.")
}
