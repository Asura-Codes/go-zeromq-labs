package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"gemini-zeromq-labs/lab15/internal/config"
	"github.com/go-zeromq/zmq4"
)

func dialWithRetry(socket zmq4.Socket, addr string, name string) {
	for {
		err := socket.Dial(addr)
		if err == nil {
			log.Printf("Connected to %s at %s", name, addr)
			return
		}
		log.Printf("Failed to dial %s (%s): %v. Retrying in 1s...", name, addr, err)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	log.Printf("Starting Reduce Worker...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Pull intermediate results from map workers (REDUCE VENTILATOR)
	receiver := zmq4.NewPull(ctx)
	defer receiver.Close()
	if err := receiver.Listen(config.ReduceVentilatorAddr); err != nil {
		log.Fatalf("Failed to listen for map workers: %v", err)
	}

	// SUB for STOP signal
	control := zmq4.NewSub(ctx)
	defer control.Close()
	dialWithRetry(control, config.ControlAddr, "control")
	if err := control.SetOption(zmq4.OptionSubscribe, "STOP"); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Push final counts to sink
	sender := zmq4.NewPush(ctx)
	defer sender.Close()
	dialWithRetry(sender, config.SinkAddr, "sink")

	counts := make(map[string]int)

	// Control goroutine
	go func() {
		msg, err := control.Recv()
		if err != nil {
			return
		}
		if string(msg.Frames[0]) == "STOP" {
			log.Println("Received STOP signal.")
			cancel() // Cancel the context to break receiver.Recv()
		}
	}()

	log.Println("Waiting for data...")
	for {
		msg, err := receiver.Recv()
		if err != nil {
			// Context was cancelled or socket closed
			break
		}
		if len(msg.Frames) < 2 {
			continue
		}
		word := string(msg.Frames[0])
		count, _ := strconv.Atoi(string(msg.Frames[1]))
		counts[word] += count
	}

	log.Printf("Finalizing: Pushing %d unique words to sink...", len(counts))
	for w, c := range counts {
		sender.Send(zmq4.Msg{Frames: [][]byte{[]byte(w), []byte(strconv.Itoa(c))}})
	}
	sender.Send(zmq4.Msg{Frames: [][]byte{[]byte("FINISHED")}})
	log.Println("Reduce worker finished.")
}
