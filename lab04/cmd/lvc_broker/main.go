package main

import (
	"context"
	"encoding/hex"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab04/internal/config"

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
		logger.Info("Shutdown signal received")
		cancel()
	}()

	// 1. Backend (XSUB) - Publishers connect here
	backend := zmq4.NewXSub(ctx)
	defer backend.Close()
	if err := backend.Listen(cfg.BackendBindAddr()); err != nil {
		logger.Error("Failed to bind backend", "error", err)
		os.Exit(1)
	}

	// 2. Frontend (XPUB) - Subscribers connect here
	frontend := zmq4.NewXPub(ctx)
	defer frontend.Close()
	if err := frontend.Listen(cfg.FrontendBindAddr()); err != nil {
		logger.Error("Failed to bind frontend", "error", err)
		os.Exit(1)
	}

	logger.Info("LVC Broker started", 
		"backend", cfg.BackendBindAddr(), 
		"frontend", cfg.FrontendBindAddr())

	// Last Value Cache: Topic -> Message Frames (Slice of bytes)
	cache := make(map[string][][]byte)

	type connMsg struct {
		msg       zmq4.Msg
		isBackend bool
		err       error
	}

	msgChan := make(chan connMsg)

	// Goroutine for Backend (Publishers)
	go func() {
		for {
			msg, err := backend.Recv()
			select {
			case <-ctx.Done():
				return
			case msgChan <- connMsg{msg: msg, isBackend: true, err: err}:
			}
			if err != nil {
				// If error is fatal/context, loop might exit or just retry
				if ctx.Err() != nil {
					return
				}
			}
		}
	}()

	// Goroutine for Frontend (Subscribers)
	go func() {
		for {
			msg, err := frontend.Recv()
			select {
			case <-ctx.Done():
				return
			case msgChan <- connMsg{msg: msg, isBackend: false, err: err}:
			}
			if err != nil {
				if ctx.Err() != nil {
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case cm := <-msgChan:
			if cm.err != nil {
				// Log error but continue (unless fatal)
				// logger.Error("Recv error", "isBackend", cm.isBackend, "error", cm.err)
				continue
			}

			if cm.isBackend {
				// Message from Publisher -> Backend
				msg := cm.msg
				// Expecting [Topic, Payload] (or more frames)
				if len(msg.Frames) >= 2 {
					topic := string(msg.Frames[0])
					// Update Cache
					cachedMsg := make([][]byte, len(msg.Frames))
					for i, f := range msg.Frames {
						cachedMsg[i] = make([]byte, len(f))
						copy(cachedMsg[i], f)
					}
					cache[topic] = cachedMsg
					logger.Debug("Cached update", "topic", topic)
				}
				// Forward to Frontend (Subscribers)
				frontend.Send(msg)

			} else {
				// Message from Subscriber -> Frontend
				msg := cm.msg
				// XPUB frames: Byte 0 is flag (0x01=Sub, 0x00=Unsub)
				if len(msg.Frames) > 0 {
					frame := msg.Frames[0]
					if len(frame) > 0 {
						flag := frame[0]
						topicBytes := frame[1:]
						topic := string(topicBytes)

						if flag == 1 {
							logger.Info("New subscription detected", "topic", topic)
							// Check Cache
							if lastMsgFrames, ok := cache[topic]; ok {
								logger.Info("Sending cached value", "topic", topic)
								cachedMsg := zmq4.NewMsgFrom(lastMsgFrames...)
								frontend.Send(cachedMsg)
							}
						} else {
							logger.Debug("Unsubscription detected", "topic", topic)
						}
					} else {
						logger.Warn("Received empty frame from frontend", "hex", hex.EncodeToString(frame))
					}
				}
				// Forward subscription upstream to Publishers (via Backend)
				backend.Send(msg)
			}
		}
	}
}
