package config

const (
	// MasterPublisherAddr is for real-time deltas (PUB)
	MasterPublisherAddr = "tcp://*:5557"
	// MasterSnapshotAddr is for state requests (ROUTER)
	MasterSnapshotAddr  = "tcp://*:5558"

	NodePublisherConnect = "tcp://localhost:5557"
	NodeSnapshotConnect  = "tcp://localhost:5558"
)
