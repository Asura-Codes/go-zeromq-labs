package tracing

import (
	"time"
)

type Span struct {
	TraceID     string    `json:"trace_id"`
	SpanID      string    `json:"span_id"`
	ParentID    string    `json:"parent_id,omitempty"`
	ServiceName string    `json:"service_name"`
	Operation   string    `json:"operation"`
	StartTime   time.Time `json:"start_time"`
	Duration    string    `json:"duration"`
	Status      string    `json:"status"`
	Tags        map[string]string `json:"tags,omitempty"`
}
