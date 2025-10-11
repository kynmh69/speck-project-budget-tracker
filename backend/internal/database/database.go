package database

import (
	"fmt"
	"log"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection
func Connect(databaseURL string) error {
	var err error
	
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	
	log.Println("Database connection established")
	return nil
}

// AutoMigrate runs auto migration for all models
func AutoMigrate() error {
	log.Println("Running auto migrations...")
	
	err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.Member{},
		&models.TimeEntry{},
		&models.ProjectMember{},
		&models.Budget{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	log.Println("Migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
