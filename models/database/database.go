package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Ignore error if .env file doesn't exist
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	DB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 保持表名单数形式
		},
	})

	if err != nil {
		return err
	}

	// Auto migrate the schema
	if err := DB.AutoMigrate(
		&VisitRecord{},
		&ContentStats{},
	); err != nil {
		return err
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
