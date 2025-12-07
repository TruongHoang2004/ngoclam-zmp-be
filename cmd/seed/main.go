package main

import (
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/database"
)

func main() {
	config.InitConfig()

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("🌱 Starting seeding...")

	images, err := SeedImages(db)
	if err != nil {
		log.Fatalf("Failed to seed images: %v", err)
	}

	categories, err := SeedCategories(db, images)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	if err := SeedProducts(db, categories, images); err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	log.Println("🎉 Seeding completed successfully!")
}
