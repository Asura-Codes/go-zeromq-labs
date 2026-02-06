package config

const (
	// Titanic Broker addresses
	TitanicAddr = "tcp://127.0.0.1:5555" // Client facing
	
	// Majordomo Broker address (Titanic acts as a client to Majordomo)
	MajordomoAddr = "tcp://127.0.0.1:5556"
	
	// Storage Service address (Internal Titanic communication)
	StorageAddr = "tcp://127.0.0.1:5557"
)
