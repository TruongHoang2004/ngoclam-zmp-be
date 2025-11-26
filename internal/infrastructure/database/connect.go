package database

import (
	"fmt"
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	log.Printf("🔌 Connecting to database %s", config.AppConfig.DBUrl)
	db, err := gorm.Open(postgres.Open(config.AppConfig.DBUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("✅ Database connected")
	return db, nil
}
