package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string    `gorm:"type:varchar(255);unique" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Price       int64     `gorm:"type:bigint" json:"price,omitempty"`

	Variants      []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	ProductImages []ProductImage   `gorm:"foreignKey:ProductID" json:"product_images,omitempty"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

type ProductVariant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index" json:"product_id"`
	Name      string    `gorm:"type:varchar(255)" json:"name" `
	Stock     int64     `gorm:"type:bigint" json:"stock,omitempty"`
	Price     int64     `gorm:"type:bigint" json:"price,omitempty"`
	Order     int       `gorm:"default:0" json:"order,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	ImageID   uint      `gorm:"index;not null" json:"image_id"`
	Image     *Image    `gorm:"foreignKey:ImageID" json:"image,omitempty"`
	Order     int       `gorm:"default:0" json:"order"`
	IsMain    bool      `gorm:"default:false" json:"is_main"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
