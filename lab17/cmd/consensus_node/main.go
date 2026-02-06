package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gemini-zeromq-labs/lab17/internal/config"
	"gemini-zeromq-labs/lab17/internal/protocol"
	"github.com/go-zeromq/zmq4"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Node struct {
	ID          int
	CurrentTerm int
	VotedFor    int // Who I voted for in this term
	State       State
	
	pub      zmq4.Socket
	sub      zmq4.Socket
	
	electionTimer  *time.Timer
	heartbeatTimer *time.Ticker
	
	// Vote tracking (volatile)
	votesReceived map[int]bool // Who has voted for me?

	mu sync.Mutex
}

func main() {
	id := flag.Int("id", 0, "Node ID (1, 2, or 3). 0 for Cluster Mode.")
	flag.Parse()

	if *id == 0 {
		runCluster()
	} else {
		runNode(*id)
	}
}

func runCluster() {
	log.Println("Starting Cluster (Nodes 1, 2, 3) in Mesh Mode...")
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { runNode(1); wg.Done() }()
	go func() { runNode(2); wg.Done() }()
	go func() { runNode(3); wg.Done() }()
	wg.Wait()
}

func runNode(id int) {
	node := &Node{
		ID:            id,
		State:         Follower,
		VotedFor:      -1,
		votesReceived: make(map[int]bool),
	}

	// Initialize Sockets (PUB/SUB Mesh)
	node.initSockets()

	// Timers
	node.heartbeatTimer = time.NewTicker(config.HeartbeatInterval)
	node.heartbeatTimer.Stop() // Only runs if Leader
	
	node.electionTimer = time.NewTimer(1 * time.Hour) // Placeholder
	node.resetElectionTimer()

	log.Printf("[Node %d] Started. Term: %d. State: Follower", id, node.CurrentTerm)

	// Message Loop
	go node.listenLoop()

	// Main Loop
	for {
		select {
		case <-node.electionTimer.C:
			node.startElection()
		case <-node.heartbeatTimer.C:
			if node.State == Leader {
				node.broadcast(protocol.EventHeartbeat, 0)
			}
		}
	}
}

func (n *Node) initSockets() {
	// 1. Publisher (My Voice)
	n.pub = zmq4.NewPub(context.Background())
	// STRICT usage of 127.0.0.1 to match Dial
	addr := config.Nodes[n.ID] // e.g., tcp://127.0.0.1:5591
	if err := n.pub.Listen(addr); err != nil {
		log.Fatalf("[Node %d] Bind failed: %v", n.ID, err)
	}

	// 2. Subscriber (Ears)
	n.sub = zmq4.NewSub(context.Background())
	for peerID, peerAddr := range config.Nodes {
		if peerID == n.ID { continue }
		
		// Retry connect
		go func(a string) {
			for {
				if err := n.sub.Dial(a); err == nil { break }
				time.Sleep(200 * time.Millisecond)
			}
		}(peerAddr)
	}
	// Subscribe to everything
	n.sub.SetOption(zmq4.OptionSubscribe, "") 

	// 3. Client Gateway (PULL) - Port + 100
	// e.g., 5591 -> 5691
	gateway := zmq4.NewPull(context.Background())
	gwPort := fmt.Sprintf("tcp://127.0.0.1:%d", 5690+n.ID)
	if err := gateway.Listen(gwPort); err != nil {
		log.Printf("[Node %d] Gateway bind failed: %v", n.ID, err)
	} else {
		log.Printf("[Node %d] Gateway listening on %s", n.ID, gwPort)
		go func() {
			for {
				msg, err := gateway.Recv()
				if err != nil { continue }
				// Assume payload is just the command string
				if len(msg.Frames) > 0 {
					cmd := string(msg.Frames[0])
					log.Printf("[Node %d] Gateway received: %s. Broadcasting...", n.ID, cmd)
					
					// Broadcast to cluster (including self via SUB loop usually, but here we might need to handle self?)
					// Actually, PUB goes to peers. If I want to handle it, I should handle it.
					// But easier to just PUB it as CLIENT_CMD. My SUB will see it if I am subscribed to myself?
					// I am NOT subscribed to myself in initSockets.
					// So I should handle it locally AND broadcast.
					
					n.broadcast(protocol.EventClientCmd, 0, cmd)
					
					// Local handle (optional, if broadcast doesn't loopback)
					// n.handleEvent(...)
				}
			}
		}()
	}
}

func (n *Node) listenLoop() {
	for {
		msg, err := n.sub.Recv()
		if err != nil { continue }

		if len(msg.Frames) > 0 {
			var evt protocol.Event
			if err := json.Unmarshal(msg.Frames[0], &evt); err == nil {
				n.handleEvent(evt)
			}
		}
	}
}

func (n *Node) broadcast(evtType protocol.EventType, candidateID int, cmd ...string) {
	evt := protocol.Event{
		Type:        evtType,
		Term:        n.CurrentTerm,
		SenderID:    n.ID,
		CandidateID: candidateID,
	}
	
	if len(cmd) > 0 {
		evt.Command = cmd[0]
	}

	// If Heartbeat, CandidateID is usually LeaderID (SenderID), but we can ignore it.
	if evtType == protocol.EventRequestVote {
		evt.CandidateID = n.ID
	}

	data, _ := json.Marshal(evt)
	n.pub.Send(zmq4.Msg{Frames: [][]byte{data}})
}

func (n *Node) handleEvent(evt protocol.Event) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 1. Term Check: If we see a higher term, step down immediately.
	if evt.Term > n.CurrentTerm {
		log.Printf("[Node %d] Saw higher term %d (from Node %d). Stepping down.", n.ID, evt.Term, evt.SenderID)
		n.becomeFollower(evt.Term)
	}

	switch evt.Type {
	case protocol.EventRequestVote:
		// Logic: If term is valid and I haven't voted (or voted for him), grant vote.
		// Since we already stepped down if evt.Term > CurrentTerm, we just check equality.
		if evt.Term == n.CurrentTerm {
			if n.VotedFor == -1 || n.VotedFor == evt.CandidateID {
				n.VotedFor = evt.CandidateID
				n.resetElectionTimer() // Granting a vote resets timeout
				
				// Broadcast my vote to the world
				log.Printf("[Node %d] Voting for Node %d (Term %d)", n.ID, evt.CandidateID, n.CurrentTerm)
				
				// We release lock briefly to broadcast (avoid deadlock if Send blocks, though Pub shouldn't)
				n.mu.Unlock()
				n.broadcast(protocol.EventVoteCast, evt.CandidateID)
				n.mu.Lock()
			}
		}

	case protocol.EventVoteCast:
		// Logic: If I am a Candidate and this vote is for ME, count it.
		if n.State == Candidate && evt.Term == n.CurrentTerm && evt.CandidateID == n.ID {
			if !n.votesReceived[evt.SenderID] {
				n.votesReceived[evt.SenderID] = true
				count := len(n.votesReceived)
				// log.Printf("[Node %d] Received vote from Node %d. Total: %d", n.ID, evt.SenderID, count)

				// Quorum Check (Hardcoded for 3 nodes: 2 votes)
				if count >= 2 {
					log.Printf("[Node %d] !!! BECAME LEADER (Term %d) !!!", n.ID, n.CurrentTerm)
					n.State = Leader
					n.heartbeatTimer.Reset(config.HeartbeatInterval)
					n.electionTimer.Stop()
					
					// Immediately assert dominance
					n.mu.Unlock()
					n.broadcast(protocol.EventHeartbeat, n.ID)
					n.mu.Lock()
				}
			}
		}

	case protocol.EventHeartbeat:
		// Logic: If from current leader, reset timeout.
		if evt.Term == n.CurrentTerm {
			if n.State == Candidate {
				n.becomeFollower(evt.Term)
			}
			n.resetElectionTimer()
			// log.Printf("[Node %d] Heartbeat from Leader %d", n.ID, evt.SenderID)
		}

	case protocol.EventClientCmd:
		log.Printf("[Node %d] Received CMD: %s", n.ID, evt.Command)
		if n.State == Leader {
			log.Printf("[Node %d] I am LEADER. Executing: %s", n.ID, evt.Command)
		}
	}
}

func (n *Node) startElection() {
	n.mu.Lock()
	n.State = Candidate
	n.CurrentTerm++
	n.VotedFor = n.ID // Vote for self
	n.votesReceived = make(map[int]bool)
	n.votesReceived[n.ID] = true
	log.Printf("[Node %d] Starting Election. Term %d", n.ID, n.CurrentTerm)
	n.mu.Unlock()

	n.resetElectionTimer()
	
	// Broadcast Request
	n.broadcast(protocol.EventRequestVote, n.ID)
}

func (n *Node) becomeFollower(term int) {
	n.State = Follower
	n.CurrentTerm = term
	n.VotedFor = -1
	n.votesReceived = make(map[int]bool)
	n.heartbeatTimer.Stop()
	n.resetElectionTimer()
}

func (n *Node) resetElectionTimer() {
	if !n.electionTimer.Stop() {
		select {
		case <-n.electionTimer.C:
		default:
		}
	}
	// Random timeout 1500-3000ms
	timeout := time.Duration(1500+rand.Intn(1500)) * time.Millisecond
	n.electionTimer.Reset(timeout)
}
