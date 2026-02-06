package protocol

import (
	"encoding/json"
	"time"
)

// Metric represents a system metric message.
type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Host      string    `json:"host"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
	Status    string    `json:"status"`
}

// ToJSON serializes the metric to JSON bytes.
func (m *Metric) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON deserializes JSON bytes to a Metric.
func FromJSON(data []byte) (*Metric, error) {
	var m Metric
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
