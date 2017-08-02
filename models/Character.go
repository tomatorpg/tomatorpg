package models

import "time"

// Character stores information about characters
// (player or non-player)
type Character struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	RoomID    uint       `json:"room_id"`
	UserID    uint       `json:"user_id"`
	Name      string     `json:"name" gorm:"type:varchar(255)"`
	Desc      string     `json:"desc" gorm:"type:text"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
