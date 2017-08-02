package models

import "github.com/jinzhu/gorm"

// AutoMigrate automatically migrate database for all models
func AutoMigrate(db *gorm.DB) (err error) {
	return db.AutoMigrate(
		Room{},
		RoomActivity{},
		User{},
		UserEmail{},
		Character{},
	).Error
}
