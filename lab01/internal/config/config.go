package config

import (
	"flag"
	"os"
)

// Config holds the configuration for the application.
type Config struct {
	Endpoint string
	Topic    string
	Interval int // in seconds
}

// LoadConfig loads configuration from command-line flags or environment variables.
// Flags take precedence.
func LoadConfig() *Config {
	endpoint := flag.String("endpoint", "tcp://127.0.0.1:5555", "ZeroMQ endpoint")
	topic := flag.String("topic", "metrics", "Subscription topic")
	interval := flag.Int("interval", 2, "Publish interval in seconds")

	flag.Parse()

	// Override with env vars if set (and flags were not explicitly provided - simplified logic here, just checking env if needed)
	// For this lab, we'll let flags rule, but could add env var fallback logic.
	if envEndpoint := os.Getenv("ZMQ_ENDPOINT"); envEndpoint != "" && !isFlagPassed("endpoint") {
		*endpoint = envEndpoint
	}

	return &Config{
		Endpoint: *endpoint,
		Topic:    *topic,
		Interval: *interval,
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
