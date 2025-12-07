package main

import (
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB, categories []*model.Category, images []*model.Image) error {
	desc := "This is a sample product description."
	products := []*model.Product{
		{
			Name:        "Smartphone X",
			Description: &desc,
			Price:       999000,
			CategoryID:  categories[0].ID,
		},
		{
			Name:        "Laptop Pro",
			Description: &desc,
			Price:       1500000,
			CategoryID:  categories[0].ID,
		},
		{
			Name:        "T-Shirt Basic",
			Description: &desc,
			Price:       150000,
			CategoryID:  categories[1].ID,
		},
		{
			Name:        "Jeans Classic",
			Description: &desc,
			Price:       300000,
			CategoryID:  categories[1].ID,
		},
		{
			Name:        "Go Programming",
			Description: &desc,
			Price:       250000,
			CategoryID:  categories[2].ID,
		},
	}

	for i, p := range products {
		if err := db.FirstOrCreate(p, model.Product{Name: p.Name}).Error; err != nil {
			log.Printf("Failed to seed product %s: %v", p.Name, err)
			return err
		}

		// Seed Variants
		variants := []model.ProductVariant{
			{
				ProductID: p.ID,
				Name:      "Standard",
				Price:     p.Price,
				Stock:     100,
			},
			{
				ProductID: p.ID,
				Name:      "Premium",
				Price:     p.Price + 50000,
				Stock:     50,
			},
		}
		for _, v := range variants {
			if err := db.FirstOrCreate(&v, model.ProductVariant{ProductID: v.ProductID, Name: v.Name}).Error; err != nil {
				log.Printf("Failed to seed variant for product %s: %v", p.Name, err)
			}
		}

		// Seed Product Images
		if len(images) > i {
			productImage := model.ProductImage{
				ProductID: p.ID,
				ImageID:   images[i].ID,
				IsMain:    true,
				Order:     1,
			}
			if err := db.FirstOrCreate(&productImage, model.ProductImage{ProductID: p.ID, ImageID: images[i].ID}).Error; err != nil {
				log.Printf("Failed to seed image for product %s: %v", p.Name, err)
			}
		}
	}

	log.Println("✅ Products seeded successfully")
	return nil
}
