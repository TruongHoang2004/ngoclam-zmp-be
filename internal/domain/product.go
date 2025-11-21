package domain

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type Product struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Price       int64             `json:"price,omitempty"`
	Variants    *[]ProductVariant `json:"variants,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ProductVariant struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	Stock     int64     `json:"stock,omitempty"`
	Price     int64     `json:"price,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProductFromModel(m *model.Product) *Product {
	return &Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		// Variants will be set separately
	}
}

func (p *Product) SetVariants(variants []*model.ProductVariant) {
	var variantList []ProductVariant
	for _, v := range variants {
		variantList = append(variantList, ProductVariant{
			ID:        v.ID,
			ProductID: v.ProductID,
			Name:      v.Name,
			Stock:     v.Stock,
			Price:     v.Price,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	p.Variants = &variantList
}

func (p *Product) ToModel() *model.Product {
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (pv *ProductVariant) ToModel(productID uint) *model.ProductVariant {
	return &model.ProductVariant{
		ID:        pv.ID,
		ProductID: productID,
		Name:      pv.Name,
		Stock:     pv.Stock,
		Price:     pv.Price,
		CreatedAt: pv.CreatedAt,
		UpdatedAt: pv.UpdatedAt,
	}
}
