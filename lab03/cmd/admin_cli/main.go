package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gemini-zeromq-labs/lab03/internal/config"
	"gemini-zeromq-labs/lab03/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Get command from args
	if len(os.Args) < 2 {
		fmt.Println("Usage: admin_cli <COMMAND>")
		fmt.Println("Commands: CPU, MEM, HOST")
		os.Exit(1)
	}

	cmdStr := strings.ToUpper(os.Args[len(os.Args)-1]) // Last arg is command
	var cmdType protocol.CommandType

	switch cmdStr {
	case "CPU":
		cmdType = protocol.CMD_CPU
	case "MEM":
		cmdType = protocol.CMD_MEM
	case "HOST":
		cmdType = protocol.CMD_HOST
	default:
		fmt.Printf("Unknown command: %s. Available: CPU, MEM, HOST\n", cmdStr)
		os.Exit(1)
	}

	// Setup context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create REQ socket
	socket := zmq4.NewReq(ctx)
	defer socket.Close()

	connectAddr := cfg.ConnectAddr()
	logger.Info("Connecting to node agent", "endpoint", connectAddr)
	if err := socket.Dial(connectAddr); err != nil {
		logger.Error("Failed to connect", "error", err)
		os.Exit(1)
	}

	// Prepare Request
	req := protocol.Request{Command: cmdType}
	reqBytes, _ := req.ToBytes()

	// Send Request
	logger.Info("Sending request", "command", cmdType)
	if err := socket.Send(zmq4.NewMsg(reqBytes)); err != nil {
		logger.Error("Failed to send request", "error", err)
		os.Exit(1)
	}

	// Wait for Reply
	// Since REQ-REP is synchronous, Recv() will block until reply or context timeout
	msg, err := socket.Recv()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Error: Timeout waiting for response from Node Agent.")
		} else {
			logger.Error("Error receiving reply", "error", err)
		}
		os.Exit(1)
	}

	// Parse Response
	var resp protocol.Response
	if err := json.Unmarshal(msg.Bytes(), &resp); err != nil {
		logger.Error("Invalid response format", "error", err)
		os.Exit(1)
	}

	if resp.Status == "ERROR" {
		fmt.Printf("Remote Error: %s\n", resp.Error)
		os.Exit(1)
	}

	// Pretty print data
	printResult(cmdType, resp.Data)
}

func printResult(cmd protocol.CommandType, data interface{}) {
	fmt.Println("\n--- Remote Node Status ---")
	
	// Re-marshal to pretty print generic interface{} map
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
	fmt.Println("--------------------------")
}
