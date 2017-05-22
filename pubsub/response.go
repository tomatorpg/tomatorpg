package pubsub

// Response to RPC request
type Response struct {
	Version string      `json:"tomatorpc"`
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type"`
	Entity  string      `json:"entity"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
}
