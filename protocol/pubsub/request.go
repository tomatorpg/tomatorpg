package pubsub

import "net/url"

// Request defines the structure of an RPC call
type Request struct {
	Version string `json:"tomatorpc,omitempty"`
	ID      string `json:"id,omitempty"`
	Group   string `json:"group,omitempty"`
	Entity  string `json:"entity,omitempty"`
	Method  string `json:"method,omitempty"`

	// derived values
	Query   url.Values  `json:"-"`
	Payload interface{} `json:"-"`
}
