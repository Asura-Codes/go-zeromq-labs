package protocol

type EventType string

const (
	EventRequestVote EventType = "REQUEST_VOTE"
	EventVoteCast    EventType = "VOTE_CAST"
	EventHeartbeat   EventType = "HEARTBEAT"
	EventClientCmd   EventType = "CLIENT_CMD"
)

// Event is the single message structure for the bus
type Event struct {
	Type        EventType `json:"type"`
	Term        int       `json:"term"`
	SenderID    int       `json:"sender_id"`
	CandidateID int       `json:"candidate_id,omitempty"` // For VOTE_CAST and REQUEST_VOTE
	Command     string    `json:"command,omitempty"`      // For CLIENT_CMD
}
