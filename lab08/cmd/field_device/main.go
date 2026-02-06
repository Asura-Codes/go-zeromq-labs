package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"gemini-zeromq-labs/lab08/internal/config"
	"gemini-zeromq-labs/lab08/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

const (
	RequestTimeout = 1000 * time.Millisecond
	MaxRetries     = 3
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	deviceID := fmt.Sprintf("dev-%d", rand.Intn(1000))
	logger.Info("Starting Field Device (Reliable Client)...", "id", deviceID)

	ctx := context.Background()

	// Initial Connection
	socket := newSocket(ctx)
	defer socket.Close()

	sequence := 0

	for {
		sequence++
		req := protocol.SensorData{
			DeviceID:  deviceID,
			Value:     rand.Float64() * 100,
			Timestamp: time.Now().Unix(),
		}
		reqBytes, _ := json.Marshal(req)

		// Reliable Send Loop
		retriesLeft := MaxRetries
		success := false
		backoff := 250 * time.Millisecond

		for retriesLeft > 0 {
			// 1. Send
			// logger.Info("Sending Request", "seq", sequence, "attempt", MaxRetries-retriesLeft+1)
			msg := zmq4.NewMsg(reqBytes)
			if err := socket.Send(msg); err != nil {
				logger.Error("Send failed", "error", err)
				// Reconnect immediately if send fails
				socket.Close()
				if retriesLeft > 1 {
					time.Sleep(backoff)
					backoff *= 2
				}
				socket = newSocket(ctx)
				retriesLeft--
				continue
			}

			// 2. Poll for Reply with Timeout
			// We spawn a goroutine to read, and select on it.
			// Note: This is slightly expensive (allocating goroutine per req) but safe.
			type result struct {
				msg zmq4.Msg
				err error
			}
			resChan := make(chan result, 1) // Buffered to avoid leak if we timeout

			go func(s zmq4.Socket) {
				m, e := s.Recv()
				resChan <- result{m, e}
			}(socket)

			select {
			case res := <-resChan:
				if res.err != nil {
					logger.Error("Receive Error", "error", res.err)
					socket.Close()
					if retriesLeft > 1 {
						time.Sleep(backoff)
						backoff *= 2
					}
					socket = newSocket(ctx)
					retriesLeft--
				} else {
					// Success!
					logger.Info("Server Replied", "seq", sequence)
					success = true
					retriesLeft = 0 // Break retry loop
				}
			case <-time.After(RequestTimeout):
				logger.Warn("Timeout - Server Unresponsive", "seq", sequence)
				retriesLeft--
				// CRITICAL: Close the socket to interrupt the stuck Recv (if any) and clear state
				socket.Close()
				if retriesLeft > 0 {
					logger.Info("Retrying connection...", "backoff", backoff)
					time.Sleep(backoff)
					backoff *= 2
					socket = newSocket(ctx)
				}
			}
			
			if success {
				break
			}
		}

		if !success {
			logger.Error("Server Abandoned", "seq", sequence)
		}

		time.Sleep(1 * time.Second)
	}
}

func newSocket(ctx context.Context) zmq4.Socket {
	// REQ socket
	// We should ideally set Linger to 0 so Close() returns immediately and doesn't try to deliver pending messages
	// zmq4.WithLinger(0) ??
	// In go-zeromq, options are passed to NewReq
	socket := zmq4.NewReq(ctx)
	// if err := socket.SetOption(zmq4.OptionLinger, 0); err != nil { ... } 
	if err := socket.Dial(config.ClientConnectAddr); err != nil {
		// Log?
	}
	return socket
}
