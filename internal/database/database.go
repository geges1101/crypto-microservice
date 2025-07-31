package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграция таблиц
	if err := db.AutoMigrate(&Currency{}, &Price{}); err != nil {
		log.Printf("Error during migration: %v", err)
		return nil, err
	}

	return db, nil
}
