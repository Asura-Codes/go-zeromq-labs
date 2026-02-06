package main

import (
	"context"
	"log"

	"gemini-zeromq-labs/lab11/internal/bstar"
	"gemini-zeromq-labs/lab11/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting HA Broker [PRIMARY]...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var frontend zmq4.Socket

	// Callbacks for FSM state changes
	startService := func() {
		log.Println(">>> SERVICE STARTING (Port 5001) <<<")
		frontend = zmq4.NewRouter(ctx)
		if err := frontend.Listen(config.PrimaryFrontendAddr); err != nil {
			log.Printf("Failed to bind frontend: %v", err)
			return
		}
		
		// Start a simple echo loop for handling clients
		go func() {
			for {
				if frontend == nil { return } // Safety check
				msg, err := frontend.Recv()
				if err != nil { return } // Socket closed

				// Identity + Empty + Payload
				if len(msg.Frames) < 3 { continue }
				
				identity := msg.Frames[0]
				payload := string(msg.Frames[2])
				log.Printf("Received: %s", payload)

				// Reply
				reply := zmq4.Msg{Frames: [][]byte{identity, {}, []byte("ACK from PRIMARY")}}
				frontend.Send(reply)
			}
		}()
	}

	stopService := func() {
		log.Println(">>> SERVICE STOPPING <<<")
		if frontend != nil {
			frontend.Close()
			frontend = nil
		}
	}

	// Init Binary Star FSM
	bstar := bstar.NewBinaryStar(
		true, // Is Primary
		config.PrimaryStatePubAddr,
		config.PrimaryStateSubAddr,
	)

	// Run FSM (Blocking)
	bstar.Run(ctx, startService, stopService)
}
