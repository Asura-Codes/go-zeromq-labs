package protocol

import "time"

// LogLevel defines the severity of the log
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// LogEntry represents a raw log message to be processed
type LogEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	RawData   string    `json:"raw_data,omitempty"` // Simulating some bulky data
}

// ProcessedLogEntry represents the log after parsing/anonymization
type ProcessedLogEntry struct {
	OriginalID  string    `json:"original_id"`
	ProcessedAt time.Time `json:"processed_at"`
	Level       LogLevel  `json:"level"`
	Sanitized   bool      `json:"sanitized"`
	CleanMessage string   `json:"clean_message"`
}
