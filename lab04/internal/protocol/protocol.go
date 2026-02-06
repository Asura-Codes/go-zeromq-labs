package protocol

import (
	"encoding/json"
	"time"
)

// TelemetryData is the payload for sensor updates
type TelemetryData struct {
	SensorID  string    `json:"sensor_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Unit      string    `json:"unit"`
}

// ToBytes serializes to JSON
func (t *TelemetryData) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

// FromBytes deserializes from JSON
func FromBytes(data []byte) (*TelemetryData, error) {
	var t TelemetryData
	err := json.Unmarshal(data, &t)
	return &t, err
}
