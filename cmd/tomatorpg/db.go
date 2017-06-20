package main

import (
	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

func initDB(db *gorm.DB) {

	logger.Printf("initDB")

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserEmail{})
	db.AutoMigrate(&models.Room{})
	db.AutoMigrate(&models.RoomActivity{})

	// Create
	//db.Create(&models.User{
	//	Email: "hello+" + time.Now().Format("20060102-150405") + "@world.com",
	//})
}
