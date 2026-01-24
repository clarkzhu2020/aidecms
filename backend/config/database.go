package config

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Ensure data directory exists
	if err := os.MkdirAll("storage/database", 0755); err != nil {
		panic("failed to create database directory: " + err.Error())
	}

	// Force SQLite configuration
	dbPath := "storage/database/data.db"
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
}
