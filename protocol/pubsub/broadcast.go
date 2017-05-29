package pubsub

import "github.com/tomatorpg/tomatorpg/models"

// Broadcast is the structure for a broadcast message
type Broadcast struct {
	Version string              `json:"tomatorpc"`
	Entity  string              `json:"entity"`
	Type    string              `json:"message_type"`
	Data    models.RoomActivity `json:"data"`
}
