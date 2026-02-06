package bstar

import (
	"context"
	"log"
	"time"

	"github.com/go-zeromq/zmq4"
)

type State int

const (
	StateActive  State = iota
	StatePassive
)

func (s State) String() string {
	switch s {
	case StateActive:
		return "Active"
	case StatePassive:
		return "Passive"
	default:
		return "Unknown"
	}
}

const (
	PeerHeartbeat = "HEARTBEAT"
	PeerActive    = "I_AM_ACTIVE"
)

type BinaryStar struct {
	state          State
	pub            zmq4.Socket
	sub            zmq4.Socket
	peerAddr       string
	localAddr      string
	lastPeerActive time.Time
	isPrimary      bool
}

func NewBinaryStar(isPrimary bool, pubAddr, subAddr string) *BinaryStar {
	return &BinaryStar{
		state:     StatePassive,
		isPrimary: isPrimary,
		localAddr: pubAddr,
		peerAddr:  subAddr,
	}
}

func (bs *BinaryStar) Run(ctx context.Context, onActive func(), onPassive func()) {
	bs.pub = zmq4.NewPub(ctx)
	bs.sub = zmq4.NewSub(ctx)
	defer bs.pub.Close()
	defer bs.sub.Close()

	if err := bs.pub.Listen(bs.localAddr); err != nil {
		log.Fatalf("Failed to bind Pub: %v", err)
	}
	if err := bs.sub.Dial(bs.peerAddr); err != nil {
		log.Fatalf("Failed to dial Peer: %v", err)
	}
	if err := bs.sub.SetOption(zmq4.OptionSubscribe, ""); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Initial State
	if bs.isPrimary {
		bs.state = StateActive
		log.Println("[FSM] Role: PRIMARY -> Started as ACTIVE")
		onActive()
	} else {
		bs.state = StatePassive
		log.Println("[FSM] Role: BACKUP -> Started as PASSIVE")
		onPassive()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	bs.lastPeerActive = time.Now()

	// Handle incoming messages in a channel to select effectively
	msgChan := make(chan string)
	go func() {
		for {
			msg, err := bs.sub.Recv()
			if err != nil {
				return // Context canceled or socket closed
			}
			msgChan <- string(msg.Bytes())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-msgChan:
			bs.lastPeerActive = time.Now()
			if msg == PeerActive {
				// Peer claims to be active.
				if bs.state == StateActive {
					// Split brain or startup race.
					if !bs.isPrimary {
						// Primary wins, I back down
						log.Println("[FSM] Primary is Active. Stepping down to PASSIVE.")
						bs.state = StatePassive
						onPassive()
					} else {
						log.Println("[FSM] I am Primary and Active. Ignoring Backup's claim.")
					}
				}
			}

		case <-ticker.C:
			// Send State
			msg := PeerHeartbeat
			if bs.state == StateActive {
				msg = PeerActive
			}
			bs.pub.Send(zmq4.NewMsgString(msg))

			// Check Timeout
			if time.Since(bs.lastPeerActive) > 3500*time.Millisecond {
				if bs.state == StatePassive {
					log.Println("[FSM] Peer Timeout! Failover initiated. Becoming ACTIVE.")
					bs.state = StateActive
					onActive()
				}
			}
		}
	}
}