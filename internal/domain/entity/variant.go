package entity

import "time"

type Variant struct {
	ID        uint
	ProductID uint
	Product   Product
	SKU       string
	Price     int64

	Image     *Image
	CreatedAt time.Time
	UpdatedAt time.Time
}
