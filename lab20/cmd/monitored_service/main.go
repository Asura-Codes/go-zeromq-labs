package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gemini-zeromq-labs/lab20/internal/tracing"
	"github.com/go-zeromq/zmq4"
)

func main() {
	serviceName := flag.String("name", "worker-service", "Name of this service")
	collectorAddr := flag.String("collector", "tcp://127.0.0.1:5555", "Address of trace collector")
	flag.Parse()

	log.Printf("Service %s starting...", *serviceName)

	ctx := context.Background()
	push := zmq4.NewPush(ctx)
	defer push.Close()

	if err := push.Dial(*collectorAddr); err != nil {
		log.Fatalf("Failed to connect to collector: %v", err)
	}

	// Simulation loop
	for {
		traceID := fmt.Sprintf("tr-%d", rand.Intn(10000))
		processRequest(push, traceID, *serviceName)
		time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	}
}

func processRequest(push zmq4.Socket, traceID, serviceName string) {
	start := time.Now()
	
	// Simulate some "business logic" work
	workDuration := time.Duration(100+rand.Intn(900)) * time.Millisecond
	time.Sleep(workDuration)

	status := "OK"
	if rand.Float32() < 0.1 {
		status = "ERROR"
	}

	span := tracing.Span{
		TraceID:     traceID,
		SpanID:      fmt.Sprintf("sp-%d", rand.Intn(10000)),
		ServiceName: serviceName,
		Operation:   "handle_request",
		StartTime:   start,
		Duration:    time.Since(start).String(),
		Status:      status,
	}

	data, _ := json.Marshal(span)
	
	// PUSH the span asynchronously (side-channel)
	// In a real high-throughput app, we might use a local buffer/channel
	err := push.Send(zmq4.NewMsg(data))
	if err != nil {
		log.Printf("Failed to push telemetry: %v", err)
	}
}
