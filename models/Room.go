package models

import "github.com/jinzhu/gorm"

// Room stores information about a room
type Room struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255)"`
	ShortName string `gorm:"type:varchar(255)"`
	Password  string
}
