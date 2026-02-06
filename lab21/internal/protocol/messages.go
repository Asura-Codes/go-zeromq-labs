package protocol

type MessageType string

const (
	Acquire   MessageType = "ACQUIRE"
	Grant     MessageType = "GRANT"
	Deny      MessageType = "DENY"
	Release   MessageType = "RELEASE"
	Heartbeat MessageType = "HEARTBEAT"
)

type LockRequest struct {
	Type     MessageType `json:"type"`
	Resource string      `json:"resource"`
	ClientID string      `json:"client_id"`
	TTL      int         `json:"ttl_seconds"`
}

type LockResponse struct {
	Type     MessageType `json:"type"`
	Resource string      `json:"resource"`
	Expires  int64       `json:"expires_at,omitempty"`
}
