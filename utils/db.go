package utils

import (
	"address-book-server/logger"
	"address-book-server/model"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
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
		logger.Log.Error("Failed to connect to DB : " + err.Error(), zap.Error(err), zap.String("dsn", dsn), zap.Int64("time", time.Now().Unix()))
		panic("Failed to connect to database")
	}

	return db
}

func PerformMigration(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Address{})

	if err != nil {
		logger.Log.Error("Migration failed : " + err.Error(), zap.Error(err), zap.Int64("time", time.Now().Unix()))
		panic("Migration failed")
	}
}