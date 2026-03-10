package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"bambu-farm/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "bambu"),
		getEnv("DB_PASSWORD", "bambupassword"),
		getEnv("DB_NAME", "bambufarm"),
		getEnv("DB_PORT", "5432"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate tables for Module 2 and 3
	err = db.AutoMigrate(
		&domain.Organization{},
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.Printer{},
		&domain.PrinterStatus{},
		&domain.PrinterLog{},
		&domain.PrinterMetric{},
		&domain.PrintJob{},
		&domain.PrintHistory{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	DB = db
	return db
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
