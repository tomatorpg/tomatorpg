package models

// UserEmail contains user and email relationship
type UserEmail struct {
	ID     uint   `gorm:"primary_key"`
	UserID uint   `gorm:"index"`
	Email  string `gorm:"type:varchar(100);unique_index"`
}
