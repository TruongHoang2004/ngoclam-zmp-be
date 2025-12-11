package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CustomerInfoRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerInfo CustomerInfoRequest `json:"customer_info" binding:"required"`
	Items        []OrderItemRequest  `json:"items" binding:"required,dive"`
}

type OrderResponse struct {
	ID          uint            `json:"id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
