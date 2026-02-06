package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"gemini-zeromq-labs/lab16/internal/config"
	"github.com/go-zeromq/zmq4"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (h *Hub) run() {
	for msg := range h.broadcast {
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Websocket error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

func main() {
	log.Printf("Starting ZMQ to WebSocket Bridge...")

	hub := newHub()
	go hub.run()

	// ZMQ Subscriber
	go func() {
		sub := zmq4.NewSub(context.Background())
		defer sub.Close()

		log.Printf("Connecting to ZMQ Producer at %s", config.ZmqPubAddr)
		for {
			if err := sub.Dial(config.ZmqPubAddr); err == nil {
				break
			}
			log.Println("Waiting for ZMQ producer...")
			hub.mu.Lock() // Just a small sleep
			hub.mu.Unlock()
		}

		if err := sub.SetOption(zmq4.OptionSubscribe, "events"); err != nil {
			log.Fatalf("Failed to subscribe: %v", err)
		}

		for {
			msg, err := sub.Recv()
			if err != nil {
				log.Printf("ZMQ Recv error: %v", err)
				continue
			}
			// Frame 0 is topic, Frame 1 is payload
			if len(msg.Frames) > 1 {
				hub.broadcast <- msg.Frames[1]
			}
		}
	}()

	// WebSocket Handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}
		hub.mu.Lock()
		hub.clients[conn] = true
		hub.mu.Unlock()
		log.Printf("New WebSocket client connected. Total: %d", len(hub.clients))
	})

	// Static File Server
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("WebSocket server listening on %s", config.WsAddr)
	if err := http.ListenAndServe(config.WsAddr, nil); err != nil {
		log.Fatalf("HTTP Server error: %v", err)
	}
}
