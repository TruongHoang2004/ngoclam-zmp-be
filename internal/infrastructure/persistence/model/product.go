package model

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
)

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Price       int64     `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	CategoryID *uint          `gorm:"index"`
	Category   *Category      `gorm:"foreignKey:CategoryID"`
	Images     []ImageRelated `gorm:"polymorphic:Entity;polymorphicValue:product;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Variants   []VariantModel `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func MapProductToModel(product *entity.Product) *Product {

	var images []ImageRelated
	if product.Images != nil {
		for _, img := range product.Images {
			images = append(images, ImageRelated{
				ImageID:  img.ID,
				EntityID: product.ID,
			})
		}
	}

	var variants []VariantModel
	if product.Variants != nil {
		for _, v := range product.Variants {
			variants = append(variants, VariantModel{
				ID:        v.ID,
				SKU:       v.SKU,
				Price:     v.Price,
				ProductID: product.ID,
			})
		}
	}

	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Images:      images,
		Variants:    variants,
		CategoryID:  &product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func (p *Product) ToDomain() *entity.Product {

	var variants []entity.Variant
	for _, v := range p.Variants {
		variants = append(variants, *v.ToDomain())
	}

	var images []entity.Image
	if p.Images != nil {
		for _, imgRel := range p.Images {
			images = append(images, entity.Image{
				ID:   imgRel.ImageID,
				Path: "", // Path will be populated later
			})
		}
	}

	var category entity.Category
	if p.Category != nil {
		category = *p.Category.ToDomain()
	}

	return &entity.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category:    category,
		Images:      images,
		Variants:    variants,
	}
}
