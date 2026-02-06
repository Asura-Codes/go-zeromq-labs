package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"gemini-zeromq-labs/lab21/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

type LockEntry struct {
	ClientID  string
	ExpiresAt time.Time
}

type LockServer struct {
	locks map[string]*LockEntry
	mu    sync.Mutex
}

func main() {
	port := flag.Int("port", 5555, "Port to bind the lock server")
	flag.Parse()

	server := &LockServer{
		locks: make(map[string]*LockEntry),
	}

	ctx := context.Background()
	router := zmq4.NewRouter(ctx)
	defer router.Close()

	addr := fmt.Sprintf("tcp://*:%d", *port)
	if err := router.Listen(addr); err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	log.Printf("Lock Server listening on %s...", addr)

	// Expiration checker
	go server.purgeExpiredLocks()

	for {
		msg, err := router.Recv()
		if err != nil {
			log.Printf("Recv error: %v", err)
			continue
		}

		// ROUTER message: [Identity, Empty, Payload]
		identity := msg.Frames[0]
		payload := msg.Frames[2]

		var req protocol.LockRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			log.Printf("Malformed request: %v", err)
			continue
		}

		resp := server.handleRequest(req)
		respData, _ := json.Marshal(resp)

		err = router.Send(zmq4.NewMsgFrom(identity, []byte(""), respData))
		if err != nil {
			log.Printf("Send error: %v", err)
		}
	}
}

func (s *LockServer) handleRequest(req protocol.LockRequest) protocol.LockResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.locks[req.Resource]
	now := time.Now()

	switch req.Type {
	case protocol.Acquire:
		if exists && entry.ExpiresAt.After(now) {
			if entry.ClientID == req.ClientID {
				// Renew
				entry.ExpiresAt = now.Add(time.Duration(req.TTL) * time.Second)
				return protocol.LockResponse{Type: protocol.Grant, Resource: req.Resource, Expires: entry.ExpiresAt.Unix()}
			}
			return protocol.LockResponse{Type: protocol.Deny, Resource: req.Resource}
		}
		// Grant new lock
		expiry := now.Add(time.Duration(req.TTL) * time.Second)
		s.locks[req.Resource] = &LockEntry{ClientID: req.ClientID, ExpiresAt: expiry}
		log.Printf("LOCK GRANTED: %s to %s (TTL: %ds)", req.Resource, req.ClientID, req.TTL)
		return protocol.LockResponse{Type: protocol.Grant, Resource: req.Resource, Expires: expiry.Unix()}

	case protocol.Heartbeat:
		if exists && entry.ClientID == req.ClientID {
			entry.ExpiresAt = now.Add(time.Duration(req.TTL) * time.Second)
			return protocol.LockResponse{Type: protocol.Grant, Resource: req.Resource, Expires: entry.ExpiresAt.Unix()}
		}
		return protocol.LockResponse{Type: protocol.Deny, Resource: req.Resource}

	case protocol.Release:
		if exists && entry.ClientID == req.ClientID {
			delete(s.locks, req.Resource)
			log.Printf("LOCK RELEASED: %s by %s", req.Resource, req.ClientID)
			return protocol.LockResponse{Type: protocol.Release, Resource: req.Resource}
		}
		return protocol.LockResponse{Type: protocol.Deny, Resource: req.Resource}

	default:
		return protocol.LockResponse{Type: protocol.Deny, Resource: req.Resource}
	}
}

func (s *LockServer) purgeExpiredLocks() {
	for {
		time.Sleep(1 * time.Second)
		s.mu.Lock()
		now := time.Now()
		for res, entry := range s.locks {
			if now.After(entry.ExpiresAt) {
				log.Printf("LOCK EXPIRED: %s (held by %s)", res, entry.ClientID)
				delete(s.locks, res)
			}
		}
		s.mu.Unlock()
	}
}
