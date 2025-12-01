package dto

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
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

func NewProductVariantResponse(m *model.ProductVariant) *ProductVariantResponse {
	return &ProductVariantResponse{
		ID:        m.ID,
		ProductID: m.ProductID,
		Name:      m.Name,
		Price:     m.Price,
		// Stock:     m.Stock,
	}
}

type CreateProductRequest struct {
	Name        string                        `json:"name" binding:"required,min=1,max=255"`
	Description string                        `json:"description"`
	Price       int64                         `json:"price" binding:"required,gt=0"`
	CategoryID  uint                          `json:"category_id" binding:"required,gt=0"`
	Variants    []CreateProductVariantRequest `json:"variants,omitempty"`
	Images      []AttachProductImageRequest   `json:"images,omitempty"`
}

func (p *CreateProductRequest) ToModel() *model.Product {
	product := &model.Product{
		Name:        p.Name,
		Description: &p.Description,
		Price:       p.Price,
	}

	if len(p.Variants) > 0 {
		var variants []model.ProductVariant
		for _, v := range p.Variants {
			variants = append(variants, model.ProductVariant{
				Name:  v.Name,
				Price: v.Price,
				Stock: v.Stock,
			})
		}
		product.Variants = variants
	}

	if len(p.Images) > 0 {
		var images []model.ProductImage
		for _, img := range p.Images {
			images = append(images, model.ProductImage{
				ImageID: img.ImageID,
				Order:   0,
				IsMain:  img.IsMain,
			})
		}
		product.ProductImages = images
	}

	return product
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
	Images      []ProductImageResponse   `json:"images,omitempty"`
}

func NewProductResponse(m *model.Product) *ProductResponse {
	var desc string
	if m.Description != nil {
		desc = *m.Description
	}

	var variant []ProductVariantResponse
	if len(m.Variants) > 0 {
		for _, v := range m.Variants {
			variant = append(variant, *NewProductVariantResponse(&v))
		}
	}

	var images []ProductImageResponse
	if len(m.ProductImages) > 0 {
		for _, img := range m.ProductImages {
			images = append(images, *NewProductImageResponse(&img))
		}
	}

	return &ProductResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: desc,
		Price:       m.Price,
		Variants:    variant,
		Images:      images,
	}
}

type ProductImageResponse struct {
	ID        uint                    `json:"id"`
	ProductID uint                    `json:"product_id"`
	ImageID   uint                    `json:"image_id"`
	VariantID *uint                   `json:"variant_id,omitempty"`
	Order     int                     `json:"order"`
	IsMain    bool                    `json:"is_main"`
	Image     *ImageResponse          `json:"image,omitempty"`
	Variant   *ProductVariantResponse `json:"variant,omitempty"`
}

func NewProductImageResponse(img *model.ProductImage) *ProductImageResponse {
	if img == nil {
		return nil
	}

	var imageResp *ImageResponse
	if img.Image != nil {
		imageResp = NewImageResponse(img.Image)
	}

	var variantResp *ProductVariantResponse

	return &ProductImageResponse{
		ID:        img.ID,
		ProductID: img.ProductID,
		ImageID:   img.ImageID,
		Order:     img.Order,
		IsMain:    img.IsMain,
		Image:     imageResp,
		Variant:   variantResp,
	}
}

type AttachProductImageRequest struct {
	ImageID   uint  `json:"image_id" binding:"required,gt=0"`
	VariantID *uint `json:"variant_id,omitempty"`
	Order     *int  `json:"order,omitempty"`
	IsMain    bool  `json:"is_main"`
}

type UpdateProductImageRequest struct {
	VariantID *uint `json:"variant_id,omitempty"`
	Order     *int  `json:"order,omitempty"`
	IsMain    *bool `json:"is_main,omitempty"`
}
