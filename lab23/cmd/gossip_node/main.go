package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/go-zeromq/zmq4"
)

type NodeState struct {
	Name      string
	Load      int
	Status    string
	LastSeen  time.Time
}

func main() {
	pubAddr := flag.String("pub", "tcp://*:6666", "Address to bind for publishing status")
	nodeName := flag.String("name", "", "Unique name for this node")
	flag.Parse()

	if *nodeName == "" {
		hostname, _ := os.Hostname()
		*nodeName = fmt.Sprintf("node-%s-%d", hostname, os.Getpid())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	log.Printf("Starting Decentralized Gossip Node [%s]...", *nodeName)

	pub := zmq4.NewPub(ctx)
	defer pub.Close()
	if err := pub.Listen(*pubAddr); err != nil {
		log.Fatalf("Failed to bind pub: %v", err)
	}

	sub := zmq4.NewSub(ctx)
	defer sub.Close()
	if err := sub.SetOption(zmq4.OptionSubscribe, ""); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	for _, peer := range flag.Args() {
		log.Printf("[%s] Peer Discovery: Connecting to %s", *nodeName, peer)
		if err := sub.Dial(peer); err != nil {
			log.Printf("Warning: Failed to dial peer %s: %v", peer, err)
		}
	}

	var mu sync.Mutex
	cluster := make(map[string]NodeState)

	// Gossip Loop: Publish our simulated load and state
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				load := rand.Intn(100)
				status := fmt.Sprintf("LOAD:%d|T:%v", load, time.Now().Format("15:04:05"))
				msg := zmq4.NewMsgFrom([]byte(*nodeName), []byte(status))
				if err := pub.Send(msg); err != nil {
					log.Printf("[GOSSIP] Send Error: %v", err)
				}
			}
		}
	}()

	// Cluster Monitor Loop: Print the "Global View" periodically
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				names := make([]string, 0, len(cluster))
				for n := range cluster {
					// Cleanup stale nodes
					if time.Since(cluster[n].LastSeen) > 10*time.Second {
						delete(cluster, n)
						continue
					}
					names = append(names, n)
				}
				sort.Strings(names)

				fmt.Printf("\n--- [%s] Global Cluster View (%d nodes) ---\n", *nodeName, len(names)+1)
				for _, n := range names {
					fmt.Printf("  > %-15s | %s\n", n, cluster[n].Status)
				}
				fmt.Println("-------------------------------------------")
				mu.Unlock()
			}
		}
	}()

	// Listen Loop
	go func() {
		for {
			msg, err := sub.Recv()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}

			if len(msg.Frames) < 2 {
				continue
			}

			name := string(msg.Frames[0])
			status := string(msg.Frames[1])

			mu.Lock()
			cluster[name] = NodeState{
				Name:     name,
				Status:   status,
				LastSeen: time.Now(),
			}
			mu.Unlock()
		}
	}()

	<-ctx.Done()
	log.Printf("Node %s exiting.", *nodeName)
}