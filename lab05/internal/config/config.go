package config

const (
	// GatewayFrontendAddr is where clients connect (ROUTER)
	GatewayFrontendAddr = "tcp://*:5555"
	// GatewayBackendAddr is where workers connect (DEALER)
	GatewayBackendAddr = "tcp://*:5556"

	// WorkerConnectAddr is the address workers dial to reach the backend
	WorkerConnectAddr = "tcp://localhost:5556"
	// ClientConnectAddr is the address clients dial to reach the frontend
	ClientConnectAddr = "tcp://localhost:5555"
)
