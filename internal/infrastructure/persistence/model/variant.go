package model

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type VariantModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ProductID uint   `gorm:"index;not null"`
	SKU       string `gorm:"type:varchar(100);not null;unique"`
	Price     int64  `gorm:"not null"`
	ImageID   *uint  `gorm:"index;default:null"` // ImageID is nullable

	Image *ImageModel `gorm:"foreignKey:ImageID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (VariantModel) TableName() string {
	return "variants"
}

func MapVariantToModel(variant *entity.Variant) *VariantModel {
	var imageID *uint
	if variant.Image != nil {
		imageID = &variant.Image.ID
	}
	return &VariantModel{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		SKU:       variant.SKU,
		Price:     variant.Price,
		ImageID:   imageID,
		Image:     MapImageToModel(variant.Image),
		CreatedAt: variant.CreatedAt,
		UpdatedAt: variant.UpdatedAt,
	}
}

func (v *VariantModel) ToDomain() *entity.Variant {
	var image *entity.Image
	if v.Image != nil {
		image = v.Image.ToDomain()
	}
	return &entity.Variant{
		ID:        v.ID,
		ProductID: v.ProductID,
		SKU:       v.SKU,
		Price:     v.Price,
		Image:     image,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
