package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Host         string
	BackendPort  int // Pubs connect here
	FrontendPort int // Subs connect here
}

func LoadConfig() *Config {
	host := flag.String("host", "127.0.0.1", "Host Address")
	backendPort := flag.Int("backend-port", 5560, "Port for Publishers (XSUB)")
	frontendPort := flag.Int("frontend-port", 5561, "Port for Subscribers (XPUB)")
	flag.Parse()

	// Env var overrides omitted for brevity but recommended in prod

	return &Config{
		Host:         *host,
		BackendPort:  *backendPort,
		FrontendPort: *frontendPort,
	}
}

// Broker methods
func (c *Config) BackendBindAddr() string {
	return fmt.Sprintf("tcp://*:%d", c.BackendPort)
}

func (c *Config) FrontendBindAddr() string {
	return fmt.Sprintf("tcp://*:%d", c.FrontendPort)
}

// Client methods
func (c *Config) PubConnectAddr() string {
	return fmt.Sprintf("tcp://%s:%d", c.Host, c.BackendPort)
}

func (c *Config) SubConnectAddr() string {
	return fmt.Sprintf("tcp://%s:%d", c.Host, c.FrontendPort)
}
