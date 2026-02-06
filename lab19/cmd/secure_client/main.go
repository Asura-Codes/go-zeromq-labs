package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gemini-zeromq-labs/lab19/internal/security"
	zmq "github.com/pebbe/zmq4"
)

func loadKeys(filename string) (security.KeyPair, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return security.KeyPair{}, err
	}
	var keys security.KeyPair
	err = json.Unmarshal(data, &keys)
	return keys, err
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	duration := flag.Duration("duration", 30*time.Second, "Duration of the data exchange")
	flag.Parse()
	
	clientKeys, err := loadKeys("client_keys.json")
	if err != nil {
		log.Fatalf("Failed to load client_keys.json: %v", err)
	}

	serverKeys, err := loadKeys("server_keys.json")
	if err != nil {
		log.Fatalf("Failed to load server_keys.json: %v", err)
	}

	client, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatalf("Failed to create socket: %v", err)
	}
	defer client.Close()

	if err := client.SetCurveServerkey(serverKeys.Public); err != nil {
		log.Fatalf("Failed to set Server Public Key: %v", err)
	}
	if err := client.SetCurvePublickey(clientKeys.Public); err != nil {
		log.Fatalf("Failed to set Client Public Key: %v", err)
	}
	if err := client.SetCurveSecretkey(clientKeys.Secret); err != nil {
		log.Fatalf("Failed to set Client Secret Key: %v", err)
	}

	addr := "tcp://127.0.0.1:5900"
	if envAddr := os.Getenv("ZMQ_SERVER_ADDR"); envAddr != "" {
		addr = envAddr
	}

	log.Printf("Connecting to %s...", addr)
	if err := client.Connect(addr); err != nil {
		log.Fatalf("Connect failed: %v", err)
	}

	log.Printf("Starting %v secure data exchange...", *duration)
	
	start := time.Now()
	count := 0
	for time.Since(start) < *duration {
		count++
		msg := fmt.Sprintf("Secure Message #%d", count)
		
		if _, err := client.Send(msg, 0); err != nil {
			log.Printf("Send error: %v", err)
			break
		}

		resp, err := client.Recv(0)
		if err != nil {
			log.Printf("Recv error: %v", err)
			break
		}
		
		log.Printf("Server Response: %s", resp)
		time.Sleep(2 * time.Second)
	}

	log.Println("Secure exchange complete. Exiting.")
}