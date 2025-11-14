package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int64  `json:"price,omitempty"`
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func NewProductResponse(model *model.Product) *ProductResponse {
	var desc string
	if model.Description != nil {
		desc = *model.Description
	}

	return &ProductResponse{
		ID:          model.ID,
		Name:        model.Name,
		Description: desc,
		Price:       model.Price,
	}
}
