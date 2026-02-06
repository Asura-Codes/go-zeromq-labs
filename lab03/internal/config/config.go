package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port int
	Host string // For client connection
}

func LoadConfig() *Config {
	port := flag.Int("port", 5559, "Port to listen on/connect to")
	host := flag.String("host", "127.0.0.1", "Host to connect to (Client only)")
	flag.Parse()

	if envPort := os.Getenv("LAB03_PORT"); envPort != "" {
		// simple parse, ignoring error for brevity in lab
		fmt.Sscanf(envPort, "%d", port)
	}

	return &Config{
		Port: *port,
		Host: *host,
	}
}

func (c *Config) BindAddr() string {
	return fmt.Sprintf("tcp://*:%d", c.Port)
}

func (c *Config) ConnectAddr() string {
	return fmt.Sprintf("tcp://%s:%d", c.Host, c.Port)
}
