package main

import (
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

func SeedImages(db *gorm.DB) ([]*model.Image, error) {
	images := []*model.Image{
		{
			Name: "Image 1",
			URL:  "https://picsum.photos/200/300",
			Hash: "hash1",
		},
		{
			Name: "Image 2",
			URL:  "https://picsum.photos/200/300",
			Hash: "hash2",
		},
		{
			Name: "Image 3",
			URL:  "https://picsum.photos/200/300",
			Hash: "hash3",
		},
		{
			Name: "Image 4",
			URL:  "https://picsum.photos/200/300",
			Hash: "hash4",
		},
		{
			Name: "Image 5",
			URL:  "https://picsum.photos/200/300",
			Hash: "hash5",
		},
	}

	for _, img := range images {
		if err := db.FirstOrCreate(img, model.Image{Hash: img.Hash}).Error; err != nil {
			log.Printf("Failed to seed image %s: %v", img.Name, err)
			return nil, err
		}
	}

	log.Println("✅ Images seeded successfully")
	return images, nil
}
