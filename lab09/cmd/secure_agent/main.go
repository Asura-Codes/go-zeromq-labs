package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gemini-zeromq-labs/lab09/internal/config"
	"gemini-zeromq-labs/lab09/internal/protocol"
	"gemini-zeromq-labs/lab09/internal/security"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Secure Agent (App-Layer Encryption)...")

	// 1. Load Keys
	clientSec, err := security.LoadKey(config.ClientSecretKey)
	if err != nil {
		panic(err)
	}
	serverPub, err := security.LoadKey(config.ServerPublicKey)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	socket := zmq4.NewDealer(ctx)
	defer socket.Close()

	logger.Info("Connecting to C2 Server...")
	if err := socket.Dial(config.Endpoint); err != nil {
		logger.Error("Failed to connect", "error", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	seq := 0

	for range ticker.C {
		seq++
		cmd := protocol.Command{
			Name: "PING",
			Args: fmt.Sprintf("Sequence %d", seq),
		}
		b, _ := json.Marshal(cmd)
		
		// 2. Encrypt
		encMsg, err := security.Encrypt(b, clientSec, serverPub)
		if err != nil {
			logger.Error("Encryption failed", "error", err)
			continue
		}

		logger.Info("Sending Encrypted Command...", "seq", seq, "bytes", len(encMsg))

		// 3. Send
		// Dealer sends [Payload] (Encrypted)
		msg := zmq4.NewMsg(encMsg)
		if err := socket.Send(msg); err != nil {
			logger.Error("Send failed", "error", err)
			continue
		}

		// 4. Recv
		reply, err := socket.Recv()
		if err != nil {
			logger.Error("Recv failed", "error", err)
			continue
		}

		// Dealer Recv: [Payload]
		if len(reply.Frames) == 0 {
			continue
		}
		encReply := reply.Frames[0]

		// 5. Decrypt
		decReply, err := security.Decrypt(encReply, clientSec, serverPub)
		if err != nil {
			logger.Error("Decryption failed (Spoofed Server?)", "error", err)
			continue
		}

		logger.Info("Reply Decrypted", "payload", string(decReply))
	}
}