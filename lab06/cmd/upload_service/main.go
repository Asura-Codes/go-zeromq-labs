package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"gemini-zeromq-labs/lab06/internal/config"
	"gemini-zeromq-labs/lab06/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Upload Service (REQ Clients)...")

	var wg sync.WaitGroup
	files := []string{"report.pdf", "virus.exe", "image.jpg", "backup.zip", "notes.txt"}

	for _, f := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			scanFile(filename, logger)
		}(f)
	}

	wg.Wait()
	logger.Info("All uploads processed.")
}

func scanFile(filename string, logger *slog.Logger) {
	ctx := context.Background()
	socket := zmq4.NewReq(ctx)
	defer socket.Close()

	if err := socket.Dial(config.ClientConnectAddr); err != nil {
		logger.Error("Connect failed", "file", filename, "error", err)
		return
	}

	req := protocol.ScanRequest{
		Filename: filename,
		Content:  make([]byte, 1024), // Dummy content
	}
	reqBytes, _ := json.Marshal(req)

	// logger.Info("Requesting scan", "file", filename)
	if err := socket.Send(zmq4.NewMsg(reqBytes)); err != nil {
		logger.Error("Send failed", "file", filename, "error", err)
		return
	}

	msg, err := socket.Recv()
	if err != nil {
		logger.Error("Recv failed", "file", filename, "error", err)
		return
	}

	// Payload is last frame
	respBytes := msg.Frames[len(msg.Frames)-1]
	var resp protocol.ScanResponse
	json.Unmarshal(respBytes, &resp)

	level := slog.LevelInfo
	if resp.Result == "INFECTED" {
		level = slog.LevelWarn
	}
	logger.Log(context.Background(), level, "Scan Result", "file", resp.Filename, "status", resp.Result, "engine", resp.Engine)
}
