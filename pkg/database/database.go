package database

import (
	"fmt"
	"log"

	"github.com/halolight/halolight-api-go/internal/models"
	"github.com/halolight/halolight-api-go/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	logLevel := logger.Warn
	if cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("📦 Database connected successfully")

	// Auto migrate schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	log.Println("✅ Database schema migrated successfully")

	return db, nil
}
