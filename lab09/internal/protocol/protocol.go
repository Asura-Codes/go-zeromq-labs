package protocol

type Command struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}
