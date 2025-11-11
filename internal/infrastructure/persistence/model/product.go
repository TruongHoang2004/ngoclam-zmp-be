package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        *string   `gorm:"type:varchar(255)" json:"name,omitempty"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Price       *float64  `gorm:"type:decimal(10,2)" json:"price,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
