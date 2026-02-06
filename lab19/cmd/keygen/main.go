package main

import (
	"encoding/json"
	"log"
	"os"

	"gemini-zeromq-labs/lab19/internal/security"
	zmq "github.com/pebbe/zmq4"
)

func generateAndSave(filename string) (security.KeyPair, error) {
	pub, sec, err := zmq.NewCurveKeypair()
	if err != nil {
		return security.KeyPair{}, err
	}
	
	pair := security.KeyPair{Public: pub, Secret: sec}
	
	data, err := json.MarshalIndent(pair, "", "  ")
	if err != nil {
		return pair, err
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return pair, err
	}
	
	return pair, nil
}

func main() {
	serverKeys, err := generateAndSave("server_keys.json")
	if err != nil {
		log.Fatalf("Failed to save server keys: %v", err)
	}
	log.Printf("Generated server_keys.json (Public: %s)", serverKeys.Public)

	clientKeys, err := generateAndSave("client_keys.json")
	if err != nil {
		log.Fatalf("Failed to save client keys: %v", err)
	}
	log.Printf("Generated client_keys.json (Public: %s)", clientKeys.Public)
}

