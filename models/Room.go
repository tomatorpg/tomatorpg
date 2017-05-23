package models

import "time"

// Room stores information about a room
type Room struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Name      string `gorm:"type:varchar(255)"`
	ShortName string `gorm:"type:varchar(255)"`
	Password  string
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
