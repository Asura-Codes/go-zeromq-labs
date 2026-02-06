package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
)

func main() {
	addr := flag.String("addr", "tcp://localhost:5555", "Address to listen on")
	name := flag.String("name", "Worker", "Worker name")
	flag.Parse()

	log.Printf("[%s] Starting on %s...", *name, *addr)

	ctx := context.Background()
	rep := zmq4.NewRep(ctx)
	defer rep.Close()

	if err := rep.Listen(*addr); err != nil {
		log.Fatalf("[%s] Failed to listen: %v", *name, err)
	}

	for {
		msg, err := rep.Recv()
		if err != nil {
			log.Printf("[%s] Recv error: %v", *name, err)
			continue
		}

		// Expect [MDPC01][REQUEST][Service][Body]
		frames := msg.Frames
		if len(frames) < 4 {
			continue
		}

		service := string(frames[2])
		body := string(frames[3])
		log.Printf("[%s] Received request for %s: %s", *name, service, body)

		// Reply: [MDPC01][REPLY][Service][Body]
		reply := [][]byte{
			[]byte("MDPC01"),
			[]byte("\x03"), // REPLY
			[]byte(service),
			[]byte(fmt.Sprintf("[%s] ECHO: %s", *name, body)),
		}

		rep.Send(zmq4.Msg{Frames: reply})
	}
}
