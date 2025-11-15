package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);unique" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Price       int64     `gorm:"type:bigint" json:"price,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProductVariant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index" json:"product_id"`
	Name      string    `gorm:"type:varchar(255)" json:"name" `
	Stock     int64     `gorm:"type:bigint" json:"stock,omitempty"`
	Price     int64     `gorm:"type:bigint" json:"price,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
