package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"

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

func NewProductVariantResponse(model *model.ProductVariant) *ProductVariantResponse {
	return &ProductVariantResponse{
		ID:        model.ID,
		ProductID: model.ProductID,
		Name:      model.Name,
		Price:     model.Price,
		// Stock:     model.Stock,
	}
}

type CreateProductRequest struct {
	Name        string                        `json:"name" binding:"required,min=1,max=255"`
	Description string                        `json:"description"`
	Price       int64                         `json:"price" binding:"required,gt=0"`
	Variants    []CreateProductVariantRequest `json:"variants,omitempty"`
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

func NewProductResponse(model *model.Product) *ProductResponse {
	var desc string
	if model.Description != nil {
		desc = *model.Description
	}

	var variant []ProductVariantResponse
	if len(model.Variants) > 0 {
		for _, v := range model.Variants {
			variant = append(variant, *NewProductVariantResponse(&v))
		}
	}

	return &ProductResponse{
		ID:          model.ID,
		Name:        model.Name,
		Description: desc,
		Price:       model.Price,
		Variants:    variant,
	}
}
