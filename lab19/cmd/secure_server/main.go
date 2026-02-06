package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"gemini-zeromq-labs/lab19/internal/security"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	timeout := flag.Duration("timeout", 0, "Server shutdown timeout (0 for no timeout)")
	flag.Parse()
	
	data, err := os.ReadFile("server_keys.json")
	if err != nil {
		log.Fatalf("Failed to read server_keys.json: %v", err)
	}
	
	var keys security.KeyPair
	if err := json.Unmarshal(data, &keys); err != nil {
		log.Fatalf("Failed to parse keys: %v", err)
	}

	log.Printf("Server Public Key: %s", keys.Public)

	server, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatalf("Failed to create socket: %v", err)
	}
	defer server.Close()

	if err := server.SetCurveServer(1); err != nil {
		log.Fatalf("Failed to set Curve Server: %v", err)
	}
	if err := server.SetCurveSecretkey(keys.Secret); err != nil {
		log.Fatalf("Failed to set Secret Key: %v", err)
	}

	if err := server.Bind("tcp://*:5900"); err != nil {
		log.Fatalf("Bind failed: %v", err)
	}

	log.Println("Secure Server listening on :5900...")
	
	// Set a timeout to exit gracefully if configured
	if *timeout > 0 {
		go func() {
			time.Sleep(*timeout)
			log.Printf("Timeout of %v reached. Exiting.", *timeout)
			server.Close()
			os.Exit(0)
		}()
	}

	for {
		msg, err := server.Recv(0)
		if err != nil {
			log.Printf("Recv error: %v", err)
			continue
		}
		log.Printf("Received: %s", msg)
		
		if _, err := server.Send("Encrypted Hello!", 0); err != nil {
			log.Printf("Send error: %v", err)
		}
	}
}