package model

import "time"

// ProductImage represents the association between a product, an optional variant, and an image.
type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	ImageID   uint      `gorm:"index;not null" json:"image_id"`
	Order     int       `gorm:"default:0" json:"order"`
	IsMain    bool      `gorm:"default:false" json:"is_main"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
