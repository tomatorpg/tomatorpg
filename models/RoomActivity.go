package models

import (
	"encoding/json"
	"time"
)

// RoomActivity stores activities in a room
type RoomActivity struct {
	ID uint `gorm:"primary_key" json:"-"`

	// room of the activity
	Room   Room `json:"-"`
	RoomID uint `json:"room_id"`

	// User of the activity
	User   User `json:"-"`
	UserID uint `json:"user_id"`

	// CharacterID of the character in the room
	CharacterID uint `json:"character_id"`

	Action  string `json:"action"`
	Message string `json:"message"`

	// MetaJSON JSON data for the activity, any valid JSON
	MetaJSON json.RawMessage `gorm:"-" json:"meta,omitempty"`
	Meta     string          `json:"-"`

	Timestamp time.Time `json:"timestamp"`
}

// BeforeSave implements BeforeSave callback for gorm
func (activity *RoomActivity) BeforeSave() {
	activity.Meta = string(activity.MetaJSON)
}

// AfterFind implements AfterFind callback for gorm
func (activity *RoomActivity) AfterFind() {
	activity.MetaJSON = []byte(activity.Meta)
}
