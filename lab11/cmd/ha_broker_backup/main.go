package main

import (
	"context"
	"log"

	"gemini-zeromq-labs/lab11/internal/bstar"
	"gemini-zeromq-labs/lab11/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting HA Broker [BACKUP]...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var frontend zmq4.Socket

	// Callbacks for FSM state changes
	startService := func() {
		log.Println(">>> SERVICE STARTING (Port 5002) <<<")
		frontend = zmq4.NewRouter(ctx)
		if err := frontend.Listen(config.BackupFrontendAddr); err != nil {
			log.Printf("Failed to bind frontend: %v", err)
			return
		}
		
		go func() {
			for {
				if frontend == nil { return }
				msg, err := frontend.Recv()
				if err != nil { return }

				if len(msg.Frames) < 3 { continue }
				
				identity := msg.Frames[0]
				payload := string(msg.Frames[2])
				log.Printf("Received: %s", payload)

				reply := zmq4.Msg{Frames: [][]byte{identity, {}, []byte("ACK from BACKUP")}}
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
		false, // Is Backup
		config.BackupStatePubAddr,
		config.BackupStateSubAddr,
	)

	// Run FSM (Blocking)
	bstar.Run(ctx, startService, stopService)
}
