package mdp

// Majordomo Protocol Constants
// See: https://rfc.zeromq.org/spec/7/MDP/

const (
	ClientHeader = "MDPC01"
	WorkerHeader = "MDPW01"
)

// Worker Commands
const (
	CommandReady      = "\001"
	CommandRequest    = "\002"
	CommandReply      = "\003"
	CommandHeartbeat  = "\004"
	CommandDisconnect = "\005"
)

// Client Commands
const (
	ClientRequest = "\001" // Only one command in MDPC
)
