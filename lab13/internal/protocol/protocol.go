package protocol

const (
	TitanicHeader = "TITANIC"
	
	// Client -> Titanic commands
	CommandSave  = "TO_STORE" // Save a request
	CommandFetch = "TO_FETCH" // Fetch result
	CommandClose = "TO_CLOSE" // Mark as done and delete
	
	// Titanic -> Client replies
	ReplySuccess = "OK"
	ReplyPending = "PENDING"
	ReplyUnknown = "UNKNOWN"
)
