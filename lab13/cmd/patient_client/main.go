package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"gemini-zeromq-labs/lab13/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

func main() {
	titanicAddr := flag.String("titanic", "tcp://127.0.0.1:5555", "Titanic Broker address")
	flag.Parse()

	log.Printf("Starting Multi-Request Patient Client (Titanic: %s)...", *titanicAddr)

	ctx := context.Background()
	client := zmq4.NewReq(ctx)
	defer client.Close()

	if err := client.Dial(*titanicAddr); err != nil {
		log.Fatalf("Failed to dial Titanic: %v", err)
	}

	// 1. Submit multiple requests
	requestIDs := make([]string, 0)
	for i := 1; i <= 5; i++ {
		service := "echo"
		body := fmt.Sprintf("Request #%d", i)
		log.Printf("Submitting: %s", body)

		msg := [][]byte{
			[]byte(protocol.TitanicHeader),
			[]byte(protocol.CommandSave),
			[]byte(service),
			[]byte(body),
		}
		client.Send(zmq4.Msg{Frames: msg})
		res, _ := client.Recv()
		id := string(res.Frames[2])
		requestIDs = append(requestIDs, id)
		log.Printf("Submitted ID: %s", id)
	}

	// 2. Poll for all results
	pending := make(map[string]bool)
	for _, id := range requestIDs {
		pending[id] = true
	}

	for len(pending) > 0 {
		time.Sleep(2 * time.Second)
		for id := range pending {
			pollMsg := [][]byte{
				[]byte(protocol.TitanicHeader),
				[]byte(protocol.CommandFetch),
				[]byte(id),
			}
			client.Send(zmq4.Msg{Frames: pollMsg})
			res, _ := client.Recv()
			status := string(res.Frames[1])

			if status == "OK" {
				log.Printf("Result for %s: %s", id, string(res.Frames[2]))
				// Acknowledge
				closeMsg := [][]byte{
					[]byte(protocol.TitanicHeader),
					[]byte(protocol.CommandClose),
					[]byte(id),
				}
				client.Send(zmq4.Msg{Frames: closeMsg})
				client.Recv()
				delete(pending, id)
			}
		}
		if len(pending) > 0 {
			log.Printf("%d requests still pending...", len(pending))
		}
	}
	log.Println("All requests processed and closed.")
}
