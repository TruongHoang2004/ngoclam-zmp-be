package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CustomerInfo struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type OrderStatus string

const (
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusPaying     OrderStatus = "paying"
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipping   OrderStatus = "shipping"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID            string          `gorm:"primaryKey;type:varchar(255)" json:"id"`
	CustomerInfo  *CustomerInfo   `gorm:"serializer:json;type:json" json:"customer_info,omitempty"`
	TotalAmount   decimal.Decimal `gorm:"type:decimal(20,2)" json:"total_amount"`
	Status        OrderStatus     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	TransactionID *string         `gorm:"type:varchar(255)" json:"transaction_id,omitempty"`
	ZaloOrderID   *string         `gorm:"type:varchar(255)" json:"zalo_order_id,omitempty"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}

type ProductSnapshot struct {
	ProductID   uint            `json:"product_id"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description,omitempty"`
	ImageURL    string          `json:"image_url,omitempty"`
}

type OrderItem struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	OrderID         string           `gorm:"index;type:varchar(255)" json:"order_id"`
	ProductSnapshot *ProductSnapshot `gorm:"serializer:json;type:json" json:"product_snapshot"`
	Quantity        int              `json:"quantity"`
	Price           decimal.Decimal  `gorm:"type:decimal(20,2)" json:"price"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
