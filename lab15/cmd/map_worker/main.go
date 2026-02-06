package main

import (
	"context"
	"log"
	"strings"
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
	log.Printf("Starting Map Worker...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Pull tasks from ventilator
	receiver := zmq4.NewPull(ctx)
	defer receiver.Close()
	dialWithRetry(receiver, config.MapVentilatorAddr, "ventilator")

	// SUB for STOP signal
	control := zmq4.NewSub(ctx)
	defer control.Close()
	dialWithRetry(control, config.ControlAddr, "control")
	if err := control.SetOption(zmq4.OptionSubscribe, "STOP"); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Push results to reduce workers
	sender := zmq4.NewPush(ctx)
	defer sender.Close()
	dialWithRetry(sender, config.ReduceVentilatorAddr, "reducer")

	// Control goroutine
	go func() {
		msg, err := control.Recv()
		if err != nil { return }
		if string(msg.Frames[0]) == "STOP" {
			log.Println("Received STOP signal. Exiting.")
			cancel()
		}
	}()

	for {
		msg, err := receiver.Recv()
		if err != nil {
			break
		}

		text := string(msg.Frames[0])
		log.Printf("Processing chunk: %s", text)

		words := strings.Fields(strings.ToLower(text))
		for _, word := range words {
			// Clean word (simplified)
			word = strings.Trim(word, ". , ! ? ; : \" ( )")
			if word == "" { continue }
			
			// Send [word, 1] to reducer
			sender.Send(zmq4.Msg{Frames: [][]byte{[]byte(word), []byte("1")}})
		}
	}
	log.Println("Map worker finished.")
}
