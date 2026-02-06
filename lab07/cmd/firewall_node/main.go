package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"gemini-zeromq-labs/lab07/internal/config"
	"gemini-zeromq-labs/lab07/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

type NodeState struct {
	mu           sync.RWMutex
	sequence     int64
	localCache   map[string]string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Firewall Node...")

	ctx := context.Background()
	state := &NodeState{
		localCache: make(map[string]string),
	}

	// 1. Connect SUB socket first (to start buffering updates)
	subscriber := zmq4.NewSub(ctx)
	defer subscriber.Close()
	if err := subscriber.Dial(config.NodePublisherConnect); err != nil {
		logger.Error("Failed to connect to publisher", "error", err)
		os.Exit(1)
	}
	// Subscribe to everything
	if err := subscriber.SetOption(zmq4.OptionSubscribe, ""); err != nil {
		logger.Error("Failed to subscribe", "error", err)
	}

	// 2. Connect Snapshot socket (DEALER)
	snapshot := zmq4.NewDealer(ctx)
	defer snapshot.Close()
	if err := snapshot.Dial(config.NodeSnapshotConnect); err != nil {
		logger.Error("Failed to connect to snapshot", "error", err)
		os.Exit(1)
	}

	// 3. Fetch Snapshot
	logger.Info("Requesting state snapshot...")
	req := protocol.SnapshotRequest{Filter: "all"}
	rb, _ := json.Marshal(req)
	// DEALER send: [Empty, Payload]
	snapshot.Send(zmq4.NewMsgFrom([]byte{}, rb))

	for {
		msg, err := snapshot.Recv()
		if err != nil {
			break
		}
		// DEALER recv: [Empty, Payload]
		payload := msg.Frames[1]
		var update protocol.PolicyUpdate
		json.Unmarshal(payload, &update)

		if update.Key == "KTHXBAI" {
			logger.Info("Snapshot sync complete", "seq", update.Sequence)
			state.sequence = update.Sequence
			break
		}

		state.localCache[update.Key] = update.Value
		logger.Info("Snapshot item", "key", update.Key, "val", update.Value)
	}

	// 4. Process real-time updates
	logger.Info("Listening for real-time updates...")
	for {
		msg, err := subscriber.Recv()
		if err != nil {
			break
		}
		var update protocol.PolicyUpdate
		json.Unmarshal(msg.Frames[0], &update)

		if update.Sequence > state.sequence {
			state.mu.Lock()
			state.sequence = update.Sequence
			state.localCache[update.Key] = update.Value
			state.mu.Unlock()
			logger.Info("Policy applied", "key", update.Key, "val", update.Value, "seq", update.Sequence)
		} else {
			logger.Debug("Discarding old update", "seq", update.Sequence, "current", state.sequence)
		}
	}
}
