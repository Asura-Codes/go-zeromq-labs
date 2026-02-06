package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-zeromq/zmq4"
)

func main() {
	beaconPort := flag.Int("beacon-port", 9999, "UDP port for beaconing")
	servicePort := flag.Int("service-port", 5555, "TCP port for the secure service")
	deviceType := flag.String("type", "GenericDevice", "Type of the device")
	flag.Parse()

	nodeName, _ := os.Hostname()
	nodeName = fmt.Sprintf("%s-%s-%d", *deviceType, nodeName, os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	log.Printf("[%s] Starting Secure Device...", nodeName)

	server := zmq4.NewRep(ctx)
	defer server.Close()
	if err := server.Listen(fmt.Sprintf("tcp://*:%d", *servicePort)); err != nil {
		log.Fatalf("Failed to bind service: %v", err)
	}

	// Beaconing: Broadcast Identity + Type + Port
	go func() {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", *beaconPort))
		if err != nil {
			return
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return
		}
		defer conn.Close()

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		// Beacon format: NAME|TYPE|PORT
		payload := []byte(fmt.Sprintf("%s|%s|%d", nodeName, *deviceType, *servicePort))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				conn.Write(payload)
			}
		}
	}()

	// Handle status requests
	for {
		msg, err := server.Recv()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			continue
		}
		
		cmd := string(msg.Frames[0])
		var response string
		if cmd == "STATUS" {
			response = fmt.Sprintf("OK|UPTIME:%s", time.Since(time.Now()).String()) // Simple mock
		} else {
			response = "UNKNOWN_COMMAND"
		}

		server.Send(zmq4.NewMsgFrom([]byte(response)))
	}
}