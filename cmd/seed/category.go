package main

import (
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB, images []*model.Image) ([]*model.Category, error) {
	categories := []*model.Category{
		{
			Name: "Electronics",
			Slug: "electronics",
		},
		{
			Name: "Clothing",
			Slug: "clothing",
		},
		{
			Name: "Books",
			Slug: "books",
		},
	}

	for i, cat := range categories {
		if len(images) > i {
			cat.ImageID = &images[i].ID
		}
		if err := db.FirstOrCreate(cat, model.Category{Slug: cat.Slug}).Error; err != nil {
			log.Printf("Failed to seed category %s: %v", cat.Name, err)
			return nil, err
		}
	}

	log.Println("✅ Categories seeded successfully")
	return categories, nil
}
