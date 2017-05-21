package main

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/tomatorpg/tomatorpg/models"
)

func initDB(db *gorm.DB) {

	log.Printf("initDB")

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserEmail{})
	db.AutoMigrate(&models.Room{})

	// Create
	//db.Create(&models.User{
	//	Email: "hello+" + time.Now().Format("20060102-150405") + "@world.com",
	//})
}
