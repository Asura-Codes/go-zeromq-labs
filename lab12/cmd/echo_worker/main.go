package main

import (
	"context"
	"flag"
	"log"
	"time"

	"gemini-zeromq-labs/lab12/internal/config"
	"gemini-zeromq-labs/lab12/internal/mdp"
	"github.com/go-zeromq/zmq4"
)

func main() {
	name := flag.String("name", "Worker", "Worker name for logging")
	flag.Parse()

	log.Printf("[%s] Starting Echo Worker...", *name)

	ctx := context.Background()
	
	// Create DEALER socket
	worker := zmq4.NewDealer(ctx)
	defer worker.Close()

	if err := worker.Dial(config.BrokerBackendAddr); err != nil {
		log.Fatalf("[%s] Failed to dial broker: %v", *name, err)
	}

	// 1. Send READY [MDPW01][READY][Service]
	// DEALER doesn't send identity envelope, BROKER adds it.
	// So we send [Empty][MDPW01][READY][Service]
	// Wait, standard DEALER just sends content. BROKER sees [ID][Content].
	// BUT, MDP says Worker message starts with [Empty][MDPW01]...
	// If using DEALER, ZMQ doesn't automatically add Empty frame at start unless using REQ.
	// So we must manually add the Empty frame if protocol requires it.
	// MDPW Spec: Frame 1: Empty Frame.
	// Yes, we must send an empty frame first.
	
	send := func(command string, option []byte, body [][]byte) {
		frames := [][]byte{
			{}, // Empty frame
			[]byte(mdp.WorkerHeader),
			[]byte(command),
		}
		if option != nil {
			frames = append(frames, option)
		} else {
			// Usually ready has service, request/reply has client addr/empty/body
		}
		frames = append(frames, body...)
		worker.Send(zmq4.Msg{Frames: frames})
	}

	// Send READY
	send(mdp.CommandReady, nil, [][]byte{[]byte(config.ServiceNameEcho)})
	log.Printf("Sent READY for service: %s", config.ServiceNameEcho)

	// Heartbeat Loop
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Recv Loop (could be poller, but let's use select with channel if possible or simple loop)
	// ZMQ Recv is blocking. We can use a Goroutine to read.
	
	msgChan := make(chan zmq4.Msg)
	go func() {
		for {
			msg, err := worker.Recv()
			if err != nil { return }
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-ticker.C:
			// Send Heartbeat
			send(mdp.CommandHeartbeat, nil, nil)
		
		case msg := <-msgChan:
			frames := msg.Frames
			// [Empty][MDPW01][Command][...]
			if len(frames) < 3 { continue }
			// frames[0] empty
			header := string(frames[1])
			if header != mdp.WorkerHeader { continue }
			
			command := string(frames[2])
			if command == mdp.CommandRequest {
				// [Empty][Header][Command][Option/Empty][ClientAddr][Empty][Body]
				// 0: Empty
				// 1: Header
				// 2: Command
				// 3: Option (Empty)
				// 4: ClientAddr
				// 5: Empty
				// 6: Body
				
				if len(frames) < 7 { continue }
				clientAddr := frames[4]
				body := frames[6:]
				
				log.Printf("[%s] Processing request: %s", *name, string(body[0]))
				
				// Simulate work
				time.Sleep(100 * time.Millisecond)
				
				// Reply: [Empty][Header][Reply][ClientAddr][Empty][Body]
				
				frames := [][]byte{
					{},
					[]byte(mdp.WorkerHeader),
					[]byte(mdp.CommandReply),
					clientAddr,
					{},
				}
				frames = append(frames, body...)
				worker.Send(zmq4.Msg{Frames: frames})
			}
		}
	}
}
