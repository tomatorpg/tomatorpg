package models

import (
	"github.com/jinzhu/gorm"
)

// User object
type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255)"`
	PrimaryEmail  string `gorm:"type:varchar(100);unique_index"`
	VerifiedEmail bool
	Emails        []UserEmail
	Password      string `gorm:"type:varchar(255)"`
	IsAdmin       bool
}
