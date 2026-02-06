package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"gemini-zeromq-labs/lab15/internal/config"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Printf("Starting MapReduce Master...")

	ctx := context.Background()

	// 1. Ventilator for Map Workers
	mapVent := zmq4.NewPush(ctx)
	defer mapVent.Close()
	if err := mapVent.Listen(config.MapVentilatorAddr); err != nil {
		log.Fatalf("Failed to listen on map ventilator: %v", err)
	}

	// 2. Control for STOP signals
	control := zmq4.NewPub(ctx)
	defer control.Close()
	if err := control.Listen(config.ControlAddr); err != nil {
		log.Fatalf("Failed to listen on control: %v", err)
	}

	// 3. Sink for Reduce Workers
	sink := zmq4.NewPull(ctx)
	defer sink.Close()
	if err := sink.Listen(config.SinkAddr); err != nil {
		log.Fatalf("Failed to listen on sink: %v", err)
	}

	ready := flag.Bool("ready", false, "Skip wait for workers")
	flag.Parse()

	if !*ready {
		fmt.Println("Press Enter when workers are ready...")
		var input string
		fmt.Scanln(&input)
	}

	log.Println("Distributing tasks...")

	chunks := []string{
		"The quick brown fox jumps over the lazy dog",
		"ZeroMQ is a high-performance asynchronous messaging library",
		"Distributed systems are complex but powerful",
		"Go is a great language for building concurrent systems",
		"The quick ZeroMQ fox jumps over the Go dog",
	}

	for _, chunk := range chunks {
		mapVent.Send(zmq4.Msg{Frames: [][]byte{[]byte(chunk)}})
	}

	log.Println("Waiting for mapping phase to complete...")
	time.Sleep(3 * time.Second)

	log.Println("Broadcasting STOP to reducers...")
	control.Send(zmq4.Msg{Frames: [][]byte{[]byte("STOP")}})

	log.Println("Collecting results from sink (15s idleness timeout)...")
	
	results := make(map[string]int)
	finishedCount := 0
	expectedReducers := 1 

	msgChan := make(chan zmq4.Msg)
	errChan := make(chan error)

	// Receiver goroutine
	go func() {
		for {
			msg, err := sink.Recv()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}
	}()

	idleTimeout := 15 * time.Second
	timer := time.NewTimer(idleTimeout)

loop:
	for {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(idleTimeout)

		select {
		case msg := <-msgChan:
			if string(msg.Frames[0]) == "FINISHED" {
				finishedCount++
				log.Printf("Reducer %d finished.", finishedCount)
				if finishedCount >= expectedReducers {
					break loop
				}
			continue
			}

			if len(msg.Frames) < 2 { continue }

			word := string(msg.Frames[0])
			var count int
			fmt.Sscanf(string(msg.Frames[1]), "%d", &count)
			results[word] += count

		case err := <-errChan:
			log.Printf("Sink receive error: %v", err)
			break loop

		case <-timer.C:
			log.Printf("Idleness timeout (15s) reached. No more data.")
			break loop
		}
	}

	fmt.Println("\n--- FINAL WORD COUNTS ---")
	for w, c := range results {
		fmt.Printf("% -15s: %d\n", w, c)
	}

	log.Println("Master finished.")
}