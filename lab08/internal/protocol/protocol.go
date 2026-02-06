package protocol

type SensorData struct {
	DeviceID  string  `json:"device_id"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type Acknowledge struct {
	Status string `json:"status"` // "OK"
}
