package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab09/internal/config"
	"gemini-zeromq-labs/lab09/internal/protocol"
	"gemini-zeromq-labs/lab09/internal/security"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting C2 Server (App-Layer Encryption)...")

	// 1. Load Keys
	serverSec, err := security.LoadKey(config.ServerSecretKey)
	if err != nil {
		panic(err)
	}
	// For Lab 09, we only accept ONE specific client (Ironhouse 1-to-1 simulation)
	clientPub, err := security.LoadKey(config.ClientPublicKey)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// 2. Standard Socket (No Transport Security)
	socket := zmq4.NewRouter(ctx)
	if err := socket.Listen(config.Endpoint); err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}
	defer socket.Close()

	logger.Info("Listening for encrypted packets...", "endpoint", config.Endpoint)

	for {
		msg, err := socket.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Error("Receive error", "error", err)
			continue
		}

		// Router: [Identity, Empty, Payload] or [Identity, Payload] depending on sender.
		// We expect [Identity, Payload] from our custom Agent (Dealer).
		if len(msg.Frames) < 2 {
			continue
		}
		identity := msg.Frames[0]
		// If using REQ/Dealer-with-envelope logic, payload might be further back. 
		// Our Agent will send [Payload] on Dealer -> Router sees [Identity, Payload].
		encryptedPayload := msg.Frames[len(msg.Frames)-1]

		// 3. Decrypt
		payload, err := security.Decrypt(encryptedPayload, serverSec, clientPub)
		if err != nil {
			logger.Warn("Decryption Failed (Unauthorized/Corrupt)", "error", err, "client_id", fmtIdentity(identity))
			continue
		}

		// 4. Process
		var cmd protocol.Command
		if err := json.Unmarshal(payload, &cmd); err != nil {
			logger.Error("Invalid JSON", "error", err)
			continue
		}

		logger.Info("Command Received (Secure)", "cmd", cmd.Name, "args", cmd.Args)

		resp := protocol.Response{
			Result: "Executed: " + cmd.Name,
		}
		respBytes, _ := json.Marshal(resp)

		// 5. Encrypt Response
		encResp, err := security.Encrypt(respBytes, serverSec, clientPub)
		if err != nil {
			logger.Error("Encryption failed", "error", err)
			continue
		}

		// 6. Reply
		// Router sends: [Identity, Payload]
		reply := zmq4.NewMsgFrom(identity, encResp)
		if err := socket.Send(reply); err != nil {
			logger.Error("Send failed", "error", err)
		}
	}
}

func fmtIdentity(id []byte) string {
	if len(id) > 5 {
		return string(id[:5]) + "..."
	}
	return string(id)
}