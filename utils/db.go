package utils

import (
	"address-book-server/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	
	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	
	// Datasource name
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB : " + err.Error())
		panic("Failed to connect to database")
	}

	return db
}

func PerformMigration(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Address{})

	if err != nil {
		log.Fatal("Migration failed : " + err.Error())
		panic("Migration failed")
	}
}