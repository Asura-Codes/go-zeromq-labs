package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	addrs := []string{
		"tcp://127.0.0.1:5001",
		"tcp://127.0.0.1:5002",
		"tcp://127.0.0.1:5003",
	}
	for _, a := range addrs {
		fmt.Printf("%s: %d\n", a, crc32.ChecksumIEEE([]byte(a)))
	}
}
