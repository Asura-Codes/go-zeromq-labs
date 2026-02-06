package config

const (
	BrokerFrontendAddr = "tcp://127.0.0.1:5555" // For Clients
	BrokerBackendAddr  = "tcp://127.0.0.1:5556" // For Workers
	
	ServiceNameEcho    = "echo"
	ServiceNameReverse = "reverse"
)
