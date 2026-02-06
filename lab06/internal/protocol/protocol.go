package protocol

const (
	// WorkerReady is the signal sent by the worker to the broker
	WorkerReady = "\001" // Simple signal
)

// ScanRequest is sent by client
type ScanRequest struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"` // Simulating file content
}

// ScanResponse is sent by worker
type ScanResponse struct {
	Filename string `json:"filename"`
	Result   string `json:"result"` // "CLEAN", "INFECTED"
	Engine   string `json:"engine"`
}
