package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"

type CreateVariantDTO struct {
	SKU   string `json:"sku" binding:"required"`
	Price int64  `json:"price" binding:"required,gt=0"`
}

type VariantDTO struct {
	ID    uint   `json:"id"`
	SKU   string `json:"sku"`
	Price int64  `json:"price"`
}

type VariantResponse struct {
	Variants []VariantDTO `json:"variants"`
}

type CreateProductRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	CategoryID  uint               `json:"category_id" binding:"required"`
	Variants    []CreateVariantDTO `json:"variants" binding:"required,dive"`
	ImageIDs    []uint             `json:"image_ids"`
}

func (r *CreateProductRequest) ToDomain() *entity.Product {
	var variants []entity.Variant
	for _, v := range r.Variants {
		variants = append(variants, entity.Variant{
			SKU:   v.SKU,
			Price: v.Price,
		})
	}

	var images []entity.Image
	for _, id := range r.ImageIDs {
		images = append(images, entity.Image{
			ID: id,
		})
	}

	return &entity.Product{
		Name:        r.Name,
		Description: r.Description,
		CategoryID:  r.CategoryID,
		Variants:    variants,
		Images:      images,
	}
}

type UpdateProductRequest struct {
	ID          uint         `json:"id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description"`
	CategoryID  uint         `json:"category_id" binding:"required"`
	Variants    []VariantDTO `json:"variants" binding:"required,dive"`
	ImageIDs    []uint       `json:"image_ids"`
}

func (r *UpdateProductRequest) ToDomain() *entity.Product {
	var variants []entity.Variant
	for _, v := range r.Variants {
		variants = append(variants, entity.Variant{
			ID:    v.ID,
			SKU:   v.SKU,
			Price: v.Price,
		})
	}

	var images []entity.Image
	for _, id := range r.ImageIDs {
		images = append(images, entity.Image{
			ID: id,
		})
	}

	return &entity.Product{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CategoryID:  r.CategoryID,
		Variants:    variants,
		Images:      images,
	}
}

type ProductResponseDTO struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Price         int64  `json:"price"`
	CategoryID    uint   `json:"category_id"`
	OriginalPrice int64  `json:"original_price"`
	Image         string `json:"image"`
	Detail        string `json:"detail"`
}

func NewProductResponseDTO(product entity.Product) ProductResponseDTO {
	var variants []VariantDTO
	for _, v := range product.Variants {
		variants = append(variants, VariantDTO{
			ID:    v.ID,
			SKU:   v.SKU,
			Price: v.Price,
		})
	}

	return ProductResponseDTO{
		ID:            product.ID,
		CategoryID:    1,
		Name:          product.Name,
		Price:         100000,
		OriginalPrice: 100000, // or set based on your business logic
		Image:         "https://nhavuonngoclam.com/wp-content/uploads/2021/02/Cay-hong-xiem-ruot-do.jpg",
		Detail:        product.Description,
	}
}
