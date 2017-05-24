package pubsub

// Response to RPC request
type Response struct {
	Version string      `json:"tomatorpc"`
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type"`
	Entity  string      `json:"entity"`
	Action  string      `json:"action"`
	Status  string      `json:"status,omitempty"`
	Err     error       `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
