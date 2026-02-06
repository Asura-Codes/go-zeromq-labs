package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"gemini-zeromq-labs/lab05/internal/config"
	"gemini-zeromq-labs/lab05/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Audit Clients (REQ)...")

	// Simulate 5 concurrent clients
	clientCount := 5
	var wg sync.WaitGroup
	wg.Add(clientCount)

	start := time.Now()

	for i := 0; i < clientCount; i++ {
		clientID := fmt.Sprintf("client-%d", i+1)
		go func(id string) {
			defer wg.Done()
			runClient(id, logger)
		}(clientID)
	}

	wg.Wait()
	logger.Info("All clients finished", "total_duration", time.Since(start))
}

func runClient(id string, logger *slog.Logger) {
	ctx := context.Background()
	socket := zmq4.NewReq(ctx)
	defer socket.Close()

	if err := socket.Dial(config.ClientConnectAddr); err != nil {
		logger.Error("Client connect failed", "client_id", id, "error", err)
		return
	}

	// Create Audit Log
	logEntry := protocol.AuditLog{
		ID:        fmt.Sprintf("log-%s-%d", id, rand.Intn(1000)),
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   "User login attempt",
		Source:    id,
	}

	reqBytes, _ := json.Marshal(logEntry)

	logger.Info("Sending log", "client_id", id, "log_id", logEntry.ID)
	
	// Send
	if err := socket.Send(zmq4.NewMsg(reqBytes)); err != nil {
		logger.Error("Send failed", "client_id", id, "error", err)
		return
	}

	// Receive Reply
	// REQ socket handles the Envelope logic internally usually, 
	// but let's see what we get back. 
	// The ROUTER/DEALER proxy preserves the REP/REQ illusion.
	msg, err := socket.Recv()
	if err != nil {
		logger.Error("Recv failed", "client_id", id, "error", err)
		return
	}

	// Parse Response
	// Payload is the last frame
	payload := msg.Frames[len(msg.Frames)-1]
	var resp protocol.Response
	if err := json.Unmarshal(payload, &resp); err != nil {
		logger.Error("Invalid response", "client_id", id, "error", err)
		return
	}

	logger.Info("Received Ack", "client_id", id, "status", resp.Status, "worker", resp.WorkerID)
}
