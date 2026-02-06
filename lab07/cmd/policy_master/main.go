package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gemini-zeromq-labs/lab07/internal/config"
	"gemini-zeromq-labs/lab07/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

type MasterState struct {
	mu       sync.RWMutex
	sequence int64
	policies map[string]string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Policy Master...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state := &MasterState{
		policies: make(map[string]string),
		sequence: 0,
	}

	// 1. Setup Sockets
	publisher := zmq4.NewPub(ctx)
	defer publisher.Close()
	publisher.Listen(config.MasterPublisherAddr)

	snapshot := zmq4.NewRouter(ctx)
	defer snapshot.Close()
	snapshot.Listen(config.MasterSnapshotAddr)

	// 2. Snapshot Handler (ROUTER)
	go func() {
		for {
			msg, err := snapshot.Recv()
			if err != nil {
				return
			}
			// ROUTER msg: [Identity, Empty, Request]
			if len(msg.Frames) < 3 {
				continue
			}
			identity := msg.Frames[0]
			
			logger.Info("Snapshot request received", "client", string(identity))

			state.mu.RLock()
			// Send each KV as a separate message for simplicity in this lab
			// In production, you might batch them.
			for k, v := range state.policies {
				update := protocol.PolicyUpdate{
					Sequence: state.sequence,
					Key:      k,
					Value:    v,
				}
				b, _ := json.Marshal(update)
				// Send back to client: [Identity, Empty, Payload]
				snapshot.Send(zmq4.NewMsgFrom(identity, []byte{}, b))
			}
			state.mu.RUnlock()

			// Send terminator (empty key or special signal)
			terminator := protocol.PolicyUpdate{Sequence: state.sequence, Key: "KTHXBAI"}
			tb, _ := json.Marshal(terminator)
			snapshot.Send(zmq4.NewMsgFrom(identity, []byte{}, tb))
		}
	}()

	// 3. Update Loop (Simulation)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				state.mu.Lock()
				state.sequence++
				key := fmt.Sprintf("rule-%d", rand.Intn(20))
				val := fmt.Sprintf("ALLOW 10.0.0.%d", rand.Intn(255))
				state.policies[key] = val
				
				update := protocol.PolicyUpdate{
					Sequence: state.sequence,
					Key:      key,
					Value:    val,
				}
				b, _ := json.Marshal(update)
				publisher.Send(zmq4.NewMsg(b))
				
				logger.Info("Policy updated", "key", key, "val", val, "seq", state.sequence)
				state.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutting down Policy Master.")
}
