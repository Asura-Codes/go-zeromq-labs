package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"gemini-zeromq-labs/lab13/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

func main() {
	frontendAddr := flag.String("frontend", "tcp://127.0.0.1:5555", "Client-facing address")
	storageAddr := flag.String("storage", "tcp://127.0.0.1:5557", "Storage service address")
	workerList := flag.String("workers", "tcp://localhost:5560,tcp://localhost:5561", "Comma-separated list of worker addresses")
	flag.Parse()

	log.Printf("Starting Titanic Broker (Frontend: %s, Storage: %s)...", *frontendAddr, *storageAddr)

	ctx := context.Background()

	// 1. Frontend: Client-facing (ROUTER)
	frontend := zmq4.NewRouter(ctx)
	defer frontend.Close()
	if err := frontend.Listen(*frontendAddr); err != nil {
		log.Fatalf("Frontend listen failed: %v", err)
	}

	// 2. Internal: Talk to Storage Service (REQ)
	storage := zmq4.NewReq(ctx)
	defer storage.Close()
	if err := storage.Dial(*storageAddr); err != nil {
		log.Fatalf("Storage dial failed: %v", err)
	}

	// 3. Background Processing Loop
	workers := strings.Split(*workerList, ",")
	go backgroundWorker(ctx, *storageAddr, workers)

	// Frontend Handling Loop
	for {
		msg, err := frontend.Recv()
		if err != nil {
			log.Printf("Frontend recv error: %v", err)
			continue
		}

		frames := msg.Frames
		if len(frames) < 4 {
			continue
		}

		clientAddr := frames[0]
		// frames[1] empty
		header := string(frames[2])
		if header != protocol.TitanicHeader {
			continue
		}

		command := string(frames[3])
		var replyData [][]byte

		switch command {
		case protocol.CommandSave:
			// SAVE [Service] [Data...]
			if len(frames) < 6 {
				replyData = [][]byte{[]byte("ERROR"), []byte("Missing service or data")}
			} else {
				service := frames[4]
				data := frames[5]

				// Forward to Storage
				err := storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("SAVE"), service, data}})
				if err == nil {
					res, _ := storage.Recv()
					replyData = res.Frames
				} else {
					replyData = [][]byte{[]byte("ERROR")}
				}
			}

		case protocol.CommandFetch:
			// FETCH [ID]
			if len(frames) < 5 {
				replyData = [][]byte{[]byte("ERROR")}
			} else {
				id := frames[4]
				err := storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("FETCH"), id}})
				if err == nil {
					res, _ := storage.Recv()
					replyData = res.Frames
				} else {
					replyData = [][]byte{[]byte("ERROR")}
				}
			}

		case protocol.CommandClose:
			// CLOSE [ID]
			if len(frames) < 5 {
				replyData = [][]byte{[]byte("ERROR")}
			} else {
				id := frames[4]
				err := storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("DELETE"), id}})
				if err == nil {
					res, _ := storage.Recv()
					replyData = res.Frames
				} else {
					replyData = [][]byte{[]byte("ERROR")}
				}
			}
		}

		// Reply: [ClientAddr][Empty][Header][Data...]
		replyFrames := [][]byte{clientAddr, {}, []byte(protocol.TitanicHeader)}
		replyFrames = append(replyFrames, replyData...)
		frontend.Send(zmq4.Msg{Frames: replyFrames})
	}
}

func backgroundWorker(ctx context.Context, storageAddr string, workerAddrs []string) {
	log.Printf("Starting Titanic background worker with %d potential workers...", len(workerAddrs))

	// Internal storage connection for background worker
	storage := zmq4.NewReq(ctx)
	defer storage.Close()
	storage.Dial(storageAddr)

	for {
		time.Sleep(2 * time.Second)

		// 1. Get list of pending requests
		storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("LIST")}})
		res, err := storage.Recv()
		if err != nil || string(res.Frames[0]) != "OK" {
			continue
		}

		pending := res.Frames[1:]
		if len(pending) == 0 {
			continue
		}

		for _, filenameBlob := range pending {
			filename := string(filenameBlob)
			parts := strings.Split(filename, ".")
			if len(parts) < 2 {
				continue
			}
			service := parts[1]

			log.Printf("Processing pending request: %s", filename)

			// 2. Get data
			storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("GET"), []byte(filename)}})
			res, _ := storage.Recv()
			if string(res.Frames[0]) != "OK" {
				continue
			}
			data := res.Frames[1]

			// 3. Try workers
			success := false
			for _, addr := range workerAddrs {
				md := zmq4.NewReq(ctx)
				// Set a very short dial timeout
				_, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
				err := md.Dial(addr)
				cancel()
				
				if err != nil {
					md.Close()
					continue
				}

				mdMsg := [][]byte{
					[]byte("MDPC01"),
					[]byte("\x02"), // REQUEST
					[]byte(service),
					data,
				}

				md.Send(zmq4.Msg{Frames: mdMsg})
				
				// Wait for reply with timeout
				done := make(chan zmq4.Msg, 1)
				go func() {
					m, _ := md.Recv()
					done <- m
				}()

				select {
				case m := <-done:
					if len(m.Frames) >= 4 {
						replyBody := m.Frames[3]
						storage.Send(zmq4.Msg{Frames: [][]byte{[]byte("COMPLETE"), []byte(filename), replyBody}})
						storage.Recv()
						log.Printf("Request %s completed via %s", filename, addr)
						success = true
					}
				case <-time.After(1 * time.Second):
					log.Printf("Timeout on %s", addr)
				}
				md.Close()
				if success {
					break
				}
			}
		}
	}
}