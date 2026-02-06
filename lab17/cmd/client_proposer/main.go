package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/go-zeromq/zmq4"
)

func main() {
	target := flag.String("target", "tcp://127.0.0.1:5691", "Node Gateway to connect to")
	cmd := flag.String("cmd", "SET x=99", "Command to propose")
	flag.Parse()

	log.Printf("Client connecting to %s...", *target)

	// PUSH socket to send request to Gateway
	socket := zmq4.NewPush(context.Background())
	if err := socket.Dial(*target); err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	// Wait for connection
	time.Sleep(1 * time.Second)

	log.Printf("Sending command: %s", *cmd)
	if err := socket.Send(zmq4.Msg{Frames: [][]byte{[]byte(*cmd)}}); err != nil {
		log.Fatalf("Send failed: %v", err)
	}

	log.Println("Command sent to Gateway.")
	// Give it time to flush
	time.Sleep(1 * time.Second)
}
