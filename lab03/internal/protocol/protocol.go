package protocol

import "encoding/json"

// CommandType defines the type of command
type CommandType string

const (
	CMD_CPU  CommandType = "CPU"
	CMD_MEM  CommandType = "MEM"
	CMD_HOST CommandType = "HOST"
)

// Request sent by the Admin CLI
type Request struct {
	Command CommandType `json:"command"`
}

// Response sent by the Node Agent
type Response struct {
	Status string      `json:"status"` // "OK" or "ERROR"
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// CPUData represents the data payload for CMD_CPU
type CPUData struct {
	Model      string    `json:"model"`
	Cores      int       `json:"cores"`
	UsagePercent []float64 `json:"usage_percent"`
}

// MemData represents the data payload for CMD_MEM
type MemData struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
}

// HostData represents the data payload for CMD_HOST
type HostData struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Platform string `json:"platform"`
	Uptime   uint64 `json:"uptime"` // seconds
}

// Helper methods

func (r *Request) ToBytes() ([]byte, error) {
	return json.Marshal(r)
}

func FromBytes(data []byte) (*Request, error) {
	var r Request
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *Response) ToBytes() ([]byte, error) {
	return json.Marshal(r)
}
