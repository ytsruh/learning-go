package jwtauth

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	log.Println("Database is connected...")
	DB = db

	db.AutoMigrate(&User{})
}
