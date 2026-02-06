package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"gemini-zeromq-labs/lab06/internal/config"
	"gemini-zeromq-labs/lab06/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	id := "av-" + randomString(4)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("id", id)
	logger.Info("Starting AV Engine (DEALER)...")

	ctx := context.Background()
	socket := zmq4.NewDealer(ctx)
	defer socket.Close()

	if err := socket.Dial(config.WorkerConnectAddr); err != nil {
		logger.Error("Connect failed", "error", err)
		os.Exit(1)
	}

	// 1. Send Initial READY
	logger.Info("Sending READY signal...")
	readyMsg := zmq4.NewMsgFrom([]byte(protocol.WorkerReady))
	if err := socket.Send(readyMsg); err != nil {
		logger.Error("Failed to send READY", "error", err)
		return
	}

	for {
		// 2. Receive Job
		msg, err := socket.Recv()
		if err != nil {
			logger.Error("Recv failed", "error", err)
			break
		}

		// Expecting: [ClientID, Empty, RequestJSON]
		if len(msg.Frames) < 3 {
			logger.Warn("Invalid job format", "frames", len(msg.Frames))
			continue
		}

		clientID := msg.Frames[0]
		// empty := msg.Frames[1]
		reqBytes := msg.Frames[2]

		var req protocol.ScanRequest
		json.Unmarshal(reqBytes, &req)

		logger.Info("Scanning file", "filename", req.Filename, "size", len(req.Content))

		// 3. Simulate Scan
		time.Sleep(time.Duration(200+rand.Intn(800)) * time.Millisecond)
		result := "CLEAN"
		if rand.Float32() < 0.1 {
			result = "INFECTED"
		}

		resp := protocol.ScanResponse{
			Filename: req.Filename,
			Result:   result,
			Engine:   id,
		}
		respBytes, _ := json.Marshal(resp)

		// 4. Send Reply (acts as Ready)
		// Must include ClientID so Broker can route
		replyMsg := zmq4.NewMsgFrom(clientID, []byte{}, respBytes)
		if err := socket.Send(replyMsg); err != nil {
			logger.Error("Failed to send reply", "error", err)
		}
		logger.Info("Scan complete", "result", result)
	}
}

func randomString(n int) string {
	const letters = "0123456789ABCDEF"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
