package pubsub

// RPC is the structure of an RPC call
type RPC struct {
	Version string `json:"tomatorpc"`
	Context string `json:"context"`
	Entity  string `json:"entity"`
	Action  string `json:"create"`
}
