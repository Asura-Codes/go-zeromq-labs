package protocol

import "time"

// AuditLog represents the business data being processed.
type AuditLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"` // INFO, WARN, ERROR
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

// Response represents the worker's acknowledgment.
type Response struct {
	Status  string `json:"status"` // "OK", "FAILED"
	LogID   string `json:"log_id"`
	WorkerID string `json:"worker_id"`
}
