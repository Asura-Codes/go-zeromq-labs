package config

import "time"

var Nodes = map[int]string{
	1: "tcp://127.0.0.1:5591",
	2: "tcp://127.0.0.1:5592",
	3: "tcp://127.0.0.1:5593",
}

const (
	ElectionTimeoutMin = 1500 * time.Millisecond
	ElectionTimeoutMax = 3000 * time.Millisecond
	HeartbeatInterval  = 500 * time.Millisecond
)
