package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gemini-zeromq-labs/lab06/internal/config"
	"gemini-zeromq-labs/lab06/internal/protocol"

	"github.com/go-zeromq/zmq4"
)

type Event struct {
	Source string // "CLIENT" or "WORKER"
	Msg    zmq4.Msg
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting Scanner Broker (ROUTER-ROUTER)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Sockets
	frontend := zmq4.NewRouter(ctx)
	defer frontend.Close()
	frontend.Listen(config.BrokerFrontendAddr)

	backend := zmq4.NewRouter(ctx)
	defer backend.Close()
	backend.Listen(config.BrokerBackendAddr)

	// 2. Event Loop Channels
	eventChan := make(chan Event)

	// Reader Routines
	go func() {
		for {
			msg, err := frontend.Recv()
			if err != nil { return }
			select {
			case eventChan <- Event{"CLIENT", msg}:
			case <-ctx.Done(): return
			}
		}
	}()

	go func() {
		for {
			msg, err := backend.Recv()
			if err != nil { return }
			select {
			case eventChan <- Event{"WORKER", msg}:
			case <-ctx.Done(): return
			}
		}
	}()

	// 3. State
	availableWorkers := []string{} // Stack of Worker IDs (LIFO or FIFO? LRU -> FIFO)
	pendingRequests := []zmq4.Msg{}

	// Helper to dispatch
	dispatch := func() {
		for len(availableWorkers) > 0 && len(pendingRequests) > 0 {
			// Pop Worker
			workerID := availableWorkers[0]
			availableWorkers = availableWorkers[1:]

			// Pop Request
			reqMsg := pendingRequests[0]
			pendingRequests = pendingRequests[1:]

			// Construct Message to Worker: [WorkerID, ClientID, Empty, Request...]
			// reqMsg frames are [ClientID, Empty, Request...]
			// We wrap this for the Backend ROUTER which needs the destination ID as first frame.
			
			// Note: zmq4.NewMsgFrom takes variadic frames.
			// We need to construct: [WorkerID] + reqMsg.Frames
			frames := append([][]byte{[]byte(workerID)}, reqMsg.Frames...)
			msgToSend := zmq4.NewMsgFrom(frames...)

			if err := backend.Send(msgToSend); err != nil {
				logger.Error("Failed to send to worker", "worker", workerID, "error", err)
				// Put request back? Or drop? For simplicity, drop.
			} else {
				logger.Info("Dispatched job", "worker", workerID)
			}
		}
	}

	logger.Info("Broker ready.")

	// Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case evt := <-eventChan:
			if evt.Source == "WORKER" {
				// Msg: [WorkerID, Payload...]
				// Payload can be [READY] or [ClientID, Empty, Reply]
				frames := evt.Msg.Frames
				workerID := string(frames[0])
				payload := frames[1:]

				if len(payload) == 1 && string(payload[0]) == protocol.WorkerReady {
					// Worker Registration
					logger.Info("Worker Ready", "id", workerID)
					availableWorkers = append(availableWorkers, workerID)
					dispatch()
				} else if len(payload) >= 3 {
					// Reply: [ClientID, Empty, Response]
					// Route to Client
					// payload IS [ClientID, Empty, Response]
					// So just send payload as message
					replyMsg := zmq4.NewMsgFrom(payload...)
					frontend.Send(replyMsg)

					// Worker is now ready again
					availableWorkers = append(availableWorkers, workerID)
					dispatch()
				} else {
					logger.Warn("Invalid worker message", "frames", len(frames))
				}

			} else { // CLIENT
				// Msg: [ClientID, Empty, Request]
				pendingRequests = append(pendingRequests, evt.Msg)
				dispatch()
			}

		case <-sigChan:
			logger.Info("Shutting down...")
			return
		case <-ctx.Done():
			return
		}
	}
}
