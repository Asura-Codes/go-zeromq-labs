package main

import (
	"context"
	"log"

	"gemini-zeromq-labs/lab12/internal/config"
	"gemini-zeromq-labs/lab12/internal/mdp"

	"github.com/go-zeromq/zmq4"
)

// Worker represents a connected worker
type Worker struct {
	Identity string
	Service  string
}

// Service represents a named service with a list of waiting requests and available workers
type Service struct {
	Name     string
	Requests [][][]byte // List of pending requests
	Workers  []*Worker  // List of available workers
}

func main() {
	log.Println("Starting Majordomo Broker...")

	ctx := context.Background()

	// 1. Prepare Sockets
	frontend := zmq4.NewRouter(ctx)
	backend := zmq4.NewRouter(ctx)
	defer frontend.Close()
	defer backend.Close()

	if err := frontend.Listen(config.BrokerFrontendAddr); err != nil {
		log.Fatalf("Frontend bind failed: %v", err)
	}
	if err := backend.Listen(config.BrokerBackendAddr); err != nil {
		log.Fatalf("Backend bind failed: %v", err)
	}

	log.Printf("Frontend listening on %s", config.BrokerFrontendAddr)
	log.Printf("Backend listening on %s", config.BrokerBackendAddr)

	// 2. State
	services := make(map[string]*Service)
	workers := make(map[string]*Worker) // Map identity -> Worker

	// Helper to get or create service
	getService := func(name string) *Service {
		if s, ok := services[name]; ok {
			return s
		}
		s := &Service{Name: name}
		services[name] = s
		return s
	}

	// Helper to send to worker
	sendToWorker := func(worker *Worker, command string, option []byte, msg [][]byte) {
		frames := [][]byte{
			[]byte(worker.Identity),
			{},
			[]byte(mdp.WorkerHeader),
			[]byte(command),
		}
		if option != nil {
			frames = append(frames, option)
		} else {
			frames = append(frames, []byte{})
		}
		frames = append(frames, msg...)

		err := backend.Send(zmq4.Msg{Frames: frames})
		if err != nil {
			log.Printf("Error sending to worker %s: %v", worker.Identity, err)
		}
	}

	// Dispatch pending requests to workers
	dispatch := func(srv *Service, msg [][]byte) {
		if msg != nil {
			srv.Requests = append(srv.Requests, msg)
		}

		for len(srv.Workers) > 0 && len(srv.Requests) > 0 {
			worker := srv.Workers[0]
			srv.Workers = srv.Workers[1:]

			req := srv.Requests[0]
			srv.Requests = srv.Requests[1:]

			if len(req) < 2 {
				continue
			}
			clientAddr := req[0]
			body := req[2:]

			payload := append([][]byte{clientAddr, {}}, body...)
			sendToWorker(worker, mdp.CommandRequest, nil, payload)
		}
	}

	// Reactor Pattern using Channels
	type SocketMsg struct {
		Source string
		Msg    zmq4.Msg
		Err    error
	}

	msgChan := make(chan SocketMsg)

	// Frontend Reader
	go func() {
		for {
			msg, err := frontend.Recv()
			if err != nil {
				return
			}
			msgChan <- SocketMsg{Source: "FRONTEND", Msg: msg, Err: err}
		}
	}()

	// Backend Reader
	go func() {
		for {
			msg, err := backend.Recv()
			if err != nil {
				return
			}
			msgChan <- SocketMsg{Source: "BACKEND", Msg: msg, Err: err}
		}
	}()

	for {
		event := <-msgChan
		if event.Err != nil {
			continue
		}

		if event.Source == "FRONTEND" {
			// CLIENT ACTIVITY
			msg := event.Msg
			frames := msg.Frames
			if len(frames) < 6 {
				continue
			}

			clientAddr := frames[0]
			header := string(frames[2])
			if header != mdp.ClientHeader {
				log.Println("Invalid client header")
				continue
			}

			serviceName := string(frames[4])
			body := frames[5:]

			log.Printf("Client Request for Service: %s (ID: %x)", serviceName, clientAddr)

			srv := getService(serviceName)

			requestBlob := append([][]byte{clientAddr, {}}, body...)

			dispatch(srv, requestBlob)

		} else if event.Source == "BACKEND" {

			// WORKER ACTIVITY

			msg := event.Msg

			frames := msg.Frames

			// log.Printf("[DEBUG] Backend Recv: %d frames", len(frames))

			if len(frames) < 4 {

				log.Printf("[DEBUG] Backend dropped: too short")

				continue

			}

			workerAddr := string(frames[0])

			// frames[1] empty

			header := string(frames[2])

			if header != mdp.WorkerHeader {

				log.Printf("[DEBUG] Backend dropped: bad header %x", header)

				continue

			}

			command := string(frames[3])

			// Handle Commands

			switch command {

			case mdp.CommandReady:

				// [READY][Service]

				if len(frames) < 5 {
					continue
				}

				serviceName := string(frames[4])

				log.Printf("Worker %s READY for %s", workerAddr, serviceName)

				w := &Worker{Identity: workerAddr, Service: serviceName}

				workers[workerAddr] = w

				srv := getService(serviceName)

				srv.Workers = append(srv.Workers, w)

				// Trigger dispatch

				dispatch(srv, nil)

			case mdp.CommandReply:

				// log.Printf("[DEBUG] Processing REPLY from %s", workerAddr)

				// [REPLY][ClientAddr][Empty][Body]

				if len(frames) < 7 {

					log.Printf("[DEBUG] Reply too short: %d", len(frames))

					continue

				}

				clientAddr := frames[4]

				// frames[5] empty

				body := frames[6:]

				// Route to Client: [ClientAddr][Empty][MDPC01][REPLY][Service][Body]

				w, exists := workers[workerAddr]

				if !exists {

					log.Printf("[DEBUG] Unknown worker: %s", workerAddr)

					continue

				}

				log.Printf("[DEBUG] Routing to Client ID: %x", clientAddr)

				replyFrames := [][]byte{

					clientAddr,

					{},

					[]byte(mdp.ClientHeader),

					[]byte(mdp.CommandReply),

					[]byte(w.Service),
				}

				replyFrames = append(replyFrames, body...)

				err := frontend.Send(zmq4.Msg{Frames: replyFrames})

				if err != nil {

					log.Printf("[DEBUG] Frontend Send Error: %v", err)

				}

				// Worker is available again

				srv := getService(w.Service)

				srv.Workers = append(srv.Workers, w)

				dispatch(srv, nil)

			case mdp.CommandDisconnect:

				log.Printf("Worker %s Disconnected", workerAddr)
				delete(workers, workerAddr)

			case mdp.CommandHeartbeat:
				// Update liveness
			}
		}
	}
}
