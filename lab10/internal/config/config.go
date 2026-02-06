package config

import (
	"os"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Service Addresses
// defaults are set for running locally (all on localhost but different ports)
var (
	IntelPubAddress      = GetEnv("INTEL_PUB_ADDR", "tcp://*:5555")
	IntelSubAddress      = GetEnv("INTEL_SUB_ADDR", "tcp://localhost:5555")
	
	AnomalyRepAddress    = GetEnv("ANOMALY_REP_ADDR", "tcp://*:5556")
	AnomalyReqAddress    = GetEnv("ANOMALY_REQ_ADDR", "tcp://localhost:5556")
	
	AlertPullAddress     = GetEnv("ALERT_PULL_ADDR", "tcp://*:5557")
	AlertPushAddress     = GetEnv("ALERT_PUSH_ADDR", "tcp://localhost:5557")
)
