package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"gemini-zeromq-labs/lab05/internal/config"
	"gemini-zeromq-labs/lab05/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	workerID := fmtID()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("worker_id", workerID)
	logger.Info("Starting Archival Worker (DEALER)...")

	ctx := context.Background()

	// 1. Create Socket
	// DEALER socket talks to the Gateway's DEALER (Backend).
	// Actually, Gateway Backend is DEALER?
	// If Gateway Backend is DEALER, and Worker is DEALER.
	// DEALER <-> DEALER is valid.
	socket := zmq4.NewDealer(ctx)
	defer socket.Close()

	// 2. Connect to Gateway
	logger.Info("Connecting to Gateway...", "addr", config.WorkerConnectAddr)
	if err := socket.Dial(config.WorkerConnectAddr); err != nil {
		logger.Error("Failed to connect", "error", err)
		os.Exit(1)
	}

	logger.Info("Worker ready. Waiting for tasks...")

	for {
		// 3. Receive Request
		// Expecting Multipart: [Client_ID, Empty, JSON_Payload]
		msg, err := socket.Recv()
		if err != nil {
			logger.Error("Receive failed", "error", err)
			break
		}

		// Log raw frames for debugging/educational value
		// logger.Info("Received message", "frames", len(msg.Frames))

		if len(msg.Frames) < 3 {
			logger.Warn("Invalid message format, dropping", "frames", len(msg.Frames))
			continue
		}

		// The last frame is the payload
		payloadFrame := msg.Frames[len(msg.Frames)-1]
		// The frames before are the return envelope
		envelope := msg.Frames[:len(msg.Frames)-1]

		var auditLog protocol.AuditLog
		if err := json.Unmarshal(payloadFrame, &auditLog); err != nil {
			logger.Error("Failed to parse JSON", "error", err)
			continue
		}

		logger.Info("Processing Log", "log_id", auditLog.ID, "severity", auditLog.Severity)

		// 4. Simulate Work (Archival I/O)
		// Sleep 500ms - 2000ms
		sleepTime := time.Duration(500+rand.Intn(1500)) * time.Millisecond
		time.Sleep(sleepTime)

		// 5. Send Reply
		// We must send [Envelope..., Response_Payload]
		response := protocol.Response{
			Status:   "ARCHIVED",
			LogID:    auditLog.ID,
			WorkerID: workerID,
		}
		respBytes, _ := json.Marshal(response)

		// Construct reply message
		replyMsg := zmq4.NewMsgFrom(envelope...) // Copy envelope
		replyMsg.Frames = append(replyMsg.Frames, respBytes) // Append payload

		if err := socket.Send(replyMsg); err != nil {
			logger.Error("Failed to send reply", "error", err)
		}
		logger.Info("Archived", "log_id", auditLog.ID, "took", sleepTime)
	}
}

func fmtID() string {
	return "w-" + os.Getenv("COMPUTERNAME") + "-" + randomString(4)
}

func randomString(n int) string {
	const letters = "abcdef0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
