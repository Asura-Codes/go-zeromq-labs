package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"gemini-zeromq-labs/lab08/internal/config"
	"gemini-zeromq-labs/lab08/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Central Receiver (Flaky Server)...")

	ctx := context.Background()
	socket := zmq4.NewRouter(ctx)
	defer socket.Close()

	if err := socket.Listen(config.ServerAddr); err != nil {
		logger.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	for {
		msg, err := socket.Recv()
		if err != nil {
			break
		}
		
		// ROUTER: [Identity, Empty, Payload]
		if len(msg.Frames) < 3 {
			continue
		}
		identity := msg.Frames[0]
		payload := msg.Frames[2]

		var data protocol.SensorData
		json.Unmarshal(payload, &data)

		// Simulate Failures
		r := rand.Intn(100)
		if r < 20 {
			logger.Warn("SIMULATING CRASH (Dropping request)", "device", data.DeviceID)
			continue // No reply
		}
		if r < 40 {
			logger.Warn("SIMULATING OVERLOAD (Slow response)", "device", data.DeviceID)
			time.Sleep(2 * time.Second) // Assuming client timeout is 1s
		}

		// Success
		logger.Info("Received Data", "device", data.DeviceID, "val", data.Value)

		ack := protocol.Acknowledge{Status: "OK"}
		b, _ := json.Marshal(ack)
		
		// Reply
		reply := zmq4.NewMsgFrom(identity, []byte{}, b)
		socket.Send(reply)
	}
}
