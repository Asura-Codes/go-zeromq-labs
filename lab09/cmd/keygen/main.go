package main

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

func main() {
	fmt.Println("Generating Curve25519 Keypairs for Ironhouse Pattern...")
	fmt.Println("-------------------------------------------------------")

	// Generate Server Pair
	pubS, secS, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	
	// Generate Client Pair
	pubC, secC, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Println("Instructions:")
	fmt.Println("1. Copy the keys below.")
	fmt.Println("2. Paste them into 'internal/config/config.go'.")
	fmt.Println("-------------------------------------------------------")
	
	fmt.Printf("ServerPublicKey = \"%x\"\n", *pubS)
	fmt.Printf("ServerSecretKey = \"%x\"\n", *secS)
	fmt.Println()
	fmt.Printf("ClientPublicKey = \"%x\"\n", *pubC)
	fmt.Printf("ClientSecretKey = \"%x\"\n", *secC)
	fmt.Println("-------------------------------------------------------")
}

