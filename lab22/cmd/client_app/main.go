package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-zeromq/zmq4"
)

func main() {
	brokerAddr := flag.String("broker", "tcp://localhost:5555", "Address of the broker frontend")
	identity := flag.String("id", "", "Identity of the client")
	requests := flag.Int("n", 1000, "Number of requests to send")
	flag.Parse()

	if *identity == "" {
		hostname, _ := os.Hostname()
		*identity = fmt.Sprintf("client-%s-%d", hostname, os.Getpid())
	}

	ctx := context.Background()
	client := zmq4.NewReq(ctx, zmq4.WithID(zmq4.SocketIdentity(*identity)))
	defer client.Close()

	log.Printf("Client [%s] connecting to %s", *identity, *brokerAddr)
	if err := client.Dial(*brokerAddr); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	start := time.Now()
	successCount := 0

	for i := 1; i <= *requests; i++ {
		payload := fmt.Sprintf("REQ %d", i)
		if err := client.Send(zmq4.NewMsgFrom([]byte(payload))); err != nil {
			log.Printf("Send error: %v", err)
			break
		}

		_, err := client.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
		} else {
			successCount++
		}
	}

	elapsed := time.Since(start)
	log.Printf("Client [%s] finished. Successful: %d/%d in %v (%.2f req/s)", 
		*identity, successCount, *requests, elapsed, float64(successCount)/elapsed.Seconds())
}