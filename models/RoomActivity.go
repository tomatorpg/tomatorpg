package models

import "time"

// RoomActivity stores activities in a room
type RoomActivity struct {
	ID uint `gorm:"primary_key" json:"-"`

	// room of the activity
	Room   Room
	RoomID uint `json:"room_id"`

	// User of the activity
	User   User
	UserID uint `json:"user_id"`

	Action  string `json:"action"`
	Message string `json:"message"`

	Timestamp time.Time `json:"timestamp"`
}
