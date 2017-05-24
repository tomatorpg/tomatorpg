package pubsub

// Request defines the structure of an RPC call
type Request struct {
	Version string `json:"tomatorpc,omitempty"`
	ID      string `json:"id,omitempty"`
	Context string `json:"context,omitempty"`
	Entity  string `json:"entity,omitempty"`
	Action  string `json:"action,omitempty"`
}
