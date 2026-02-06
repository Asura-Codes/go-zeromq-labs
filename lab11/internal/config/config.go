package config

const (
	// Client facing ports
	PrimaryFrontendAddr = "tcp://127.0.0.1:5001"
	BackupFrontendAddr  = "tcp://127.0.0.1:5002"

	// Peering ports (State synchronization)
	// Primary binds to PrimaryStateAddr, Backup connects to it
	// Backup binds to BackupStateAddr, Primary connects to it
	// Actually, for Binary Star, usually they have a dedicated pair of sockets.
	// Let's use Pub-Sub for heartbeats/state updates for simplicity in this lab.
	
	PrimaryStatePubAddr = "tcp://127.0.0.1:6001"
	PrimaryStateSubAddr = "tcp://127.0.0.1:6002" // Backup publishes here

	BackupStatePubAddr  = "tcp://127.0.0.1:6002"
	BackupStateSubAddr  = "tcp://127.0.0.1:6001" // Primary publishes here
	
	// Heartbeat intervals
	HeartbeatInterval = "1s"
	FailoverTimeout   = "3s"
)
