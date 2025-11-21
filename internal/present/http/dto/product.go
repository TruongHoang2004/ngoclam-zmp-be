package dto

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
)

type CreateProductVariantRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Price int64  `json:"price" binding:"required,gt=0"`
	Stock int64  `json:"stock"`
}

type UpdateProductVariantRequest struct {
	ID    uint    `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Price *int64  `json:"price,omitempty"`
	Stock *int64  `json:"stock,omitempty"`
}

type AddProductVariantRequest struct {
	ProductID uint   `json:"product_id" binding:"required,gt=0"`
	Name      string `json:"name" binding:"required,min=1,max=255"`
	Price     int64  `json:"price" binding:"required,gt=0"`
	Stock     int64  `json:"stock"`
}

type ProductVariantResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	// Stock     int64  `json:"stock"`
}

func NewProductVariantResponse(domain *domain.ProductVariant) *ProductVariantResponse {
	return &ProductVariantResponse{
		ID:        domain.ID,
		ProductID: domain.ProductID,
		Name:      domain.Name,
		Price:     domain.Price,
		// Stock:     domain.Stock,
	}
}

type CreateProductRequest struct {
	Name        string                        `json:"name" binding:"required,min=1,max=255"`
	Description string                        `json:"description"`
	Price       int64                         `json:"price" binding:"required,gt=0"`
	Variants    []CreateProductVariantRequest `json:"variants,omitempty"`
}

func (p *CreateProductRequest) ToDomain() *domain.Product {
	domainProduct := &domain.Product{
		Name:        p.Name,
		Description: &p.Description,
		Price:       p.Price,
	}

	if len(p.Variants) > 0 {
		var variants []domain.ProductVariant
		for _, v := range p.Variants {
			variants = append(variants, domain.ProductVariant{
				Name:  v.Name,
				Price: v.Price,
				Stock: v.Stock,
			})
		}
		domainProduct.Variants = &variants
	}

	return domainProduct
}

type UpdateProductRequest struct {
	ID          uint    `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int64  `json:"price,omitempty"`
}

type ProductResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Price       int64                    `json:"price"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
}

func NewProductResponse(domain *domain.Product) *ProductResponse {
	var desc string
	if domain.Description != nil {
		desc = *domain.Description
	}

	var variant []ProductVariantResponse
	if (domain.Variants != nil) && len(*domain.Variants) > 0 {
		if len(*domain.Variants) > 0 {
			for _, v := range *domain.Variants {
				variant = append(variant, *NewProductVariantResponse(&v))
			}
		}
	}

	return &ProductResponse{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: desc,
		Price:       domain.Price,
		Variants:    variant,
	}
}
