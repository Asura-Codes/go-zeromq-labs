package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds the configuration for the application.
type Config struct {
	Host          string
	CollectorPort int
	SinkPort      int
}

// LoadConfig loads configuration from command-line flags.
func LoadConfig() *Config {
	host := flag.String("host", "127.0.0.1", "Host IP for workers to connect to")
	collectorPort := flag.Int("collector-port", 5557, "Port for the Log Collector (Ventilator)")
	sinkPort := flag.Int("sink-port", 5558, "Port for the Storage Writer (Sink)")

	flag.Parse()

	// Simple Env override check (optional but good for Docker)
	if v := os.Getenv("LAB02_HOST"); v != "" {
		*host = v
	}
	if v := os.Getenv("LAB02_COLLECTOR_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			*collectorPort = p
		}
	}
	if v := os.Getenv("LAB02_SINK_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			*sinkPort = p
		}
	}

	return &Config{
		Host:          *host,
		CollectorPort: *collectorPort,
		SinkPort:      *sinkPort,
	}
}

// CollectorBindAddr returns the endpoint for the Collector to bind to (listening on all interfaces).
func (c *Config) CollectorBindAddr() string {
	return fmt.Sprintf("tcp://*:%d", c.CollectorPort)
}

// CollectorConnectAddr returns the endpoint for Workers to connect to.
func (c *Config) CollectorConnectAddr() string {
	return fmt.Sprintf("tcp://%s:%d", c.Host, c.CollectorPort)
}

// SinkBindAddr returns the endpoint for the Sink to bind to.
func (c *Config) SinkBindAddr() string {
	return fmt.Sprintf("tcp://*:%d", c.SinkPort)
}

// SinkConnectAddr returns the endpoint for Workers to connect to.
func (c *Config) SinkConnectAddr() string {
	return fmt.Sprintf("tcp://%s:%d", c.Host, c.SinkPort)
}
