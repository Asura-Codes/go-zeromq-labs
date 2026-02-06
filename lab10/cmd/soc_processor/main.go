package main

import (
	"context"
	"fmt"
	"log"

	"gemini-zeromq-labs/lab10/internal/config"
	"gemini-zeromq-labs/lab10/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

func main() {
	log.Println("Starting SOC Processor Hub...")

	ctx := context.Background()

	// 1. Subscribe to Intel Feed
	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	if err := sub.Dial(config.IntelSubAddress); err != nil {
		log.Fatalf("Failed to dial Intel Provider: %v", err)
	}
	if err := sub.SetOption(zmq4.OptionSubscribe, ""); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// 2. Request-Reply for Anomaly Detection
	req := zmq4.NewReq(ctx)
	defer req.Close()
	if err := req.Dial(config.AnomalyReqAddress); err != nil {
		log.Fatalf("Failed to dial Anomaly Detector: %v", err)
	}

	// 3. Push Alerts to Logger
	push := zmq4.NewPush(ctx)
	defer push.Close()
	if err := push.Dial(config.AlertPushAddress); err != nil {
		log.Fatalf("Failed to dial Alert Logger: %v", err)
	}

	log.Println("SOC Processor connected and operational.")

	for {
		// Wait for data from intel feed
		msg, err := sub.Recv()
		if err != nil {
			log.Printf("Error receiving from sub: %v", err)
			continue
		}

		ip := string(msg.Bytes())
		log.Printf("Processing Intel: IP %s", ip)

		// Synchronously query detector
		if err := req.Send(zmq4.NewMsgString(ip)); err != nil {
			log.Printf("Error sending to detector: %v", err)
			continue
		}

		reply, err := req.Recv()
		if err != nil {
			log.Printf("Error receiving from detector: %v", err)
			continue
		}

		verdict := string(reply.Bytes())
		log.Printf("Verdict for %s: %s", ip, verdict)

		// If malicious, push to logger
		if verdict == protocol.StatusMalicious {
			alertMsg := fmt.Sprintf("Threat Detected! Source IP: %s", ip)
			if err := push.Send(zmq4.NewMsgString(alertMsg)); err != nil {
				log.Printf("Error pushing alert: %v", err)
			}
		}
	}
}
