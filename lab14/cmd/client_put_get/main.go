package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
)

func main() {
	port := flag.Int("port", 5001, "Port of a node to connect to")
	flag.Parse()

	addr := fmt.Sprintf("tcp://127.0.0.1:%d", *port)
	log.Printf("Connecting to DHT via %s...", addr)

	ctx := context.Background()
	client := zmq4.NewReq(ctx)
	defer client.Close()

	if err := client.Dial(addr); err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	testData := map[string]string{
		"apple":      "red",
		"banana":     "yellow",
		"grape":      "purple",
		"kiwi":       "green",
		"orange":     "orange",
		"blueberry":  "blue",
		"strawberry": "red",
		"mango":      "yellow",
		"pineapple":  "yellow",
		"watermelon": "green",
	}

	// 1. PUT
	for k, v := range testData {
		log.Printf("PUT %s = %s", k, v)
		err := client.Send(zmq4.Msg{Frames: [][]byte{[]byte("PUT"), []byte(k), []byte(v)}})
		if err != nil {
			log.Fatalf("Send error: %v", err)
		}

		msg, _ := client.Recv()
		log.Printf("  Response: %s", string(msg.Frames[0]))
	}

	// 2. GET
	for k, v := range testData {
		log.Printf("GET %s", k)
		client.Send(zmq4.Msg{Frames: [][]byte{[]byte("GET"), []byte(k)}})
		msg, _ := client.Recv()

		if string(msg.Frames[0]) == "OK" {
			val := string(msg.Frames[1])
			log.Printf("  Response: %s (Expected: %s)", val, v)
		} else {
			log.Printf("  Response: %s", string(msg.Frames[0]))
		}
	}
}
