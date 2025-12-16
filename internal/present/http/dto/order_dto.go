package dto

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/shopspring/decimal"
)

type CustomerInfoRequest struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Address string `json:"address" validate:"required"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type PaymentRequest struct {
	Method string `json:"method" validate:"required"`
}

type CreateOrderRequest struct {
	CustomerInfo CustomerInfoRequest `json:"customer_info" validate:"required"`
	Items        []OrderItemRequest  `json:"items" validate:"required,dive"`
	Payment      PaymentRequest      `json:"payment" validate:"required"`
}

type OrderResponse struct {
	ID          uint            `json:"id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ZaloOrderParams struct {
	Amount    int64  `json:"amount"`
	Desc      string `json:"desc"`
	Item      string `json:"item"`
	Extradata string `json:"extradata"`
	Method    string `json:"method"`
	Mac       string `json:"mac"`
}

type CreateOrderResponse struct {
	*model.Order
	ZaloParams *ZaloOrderParams `json:"zalo_params"`
	MAC        string           `json:"mac"` // Deprecated but keeping for now
}

type OrderSubmitRequest struct {
	OrderID     uint   `json:"order_id" validate:"required"`
	ZaloOrderID string `json:"zalo_order_id" validate:"required"`
}
