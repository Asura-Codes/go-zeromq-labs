package protocol

// PolicyUpdate represents a single change to the policy state.
// It is used for both Snapshot items and real-time Updates.
type PolicyUpdate struct {
	Sequence int64  `json:"sequence"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

// SnapshotRequest is sent by nodes to request the full state.
type SnapshotRequest struct {
	Filter string `json:"filter"` // e.g., which partition of policies
}
