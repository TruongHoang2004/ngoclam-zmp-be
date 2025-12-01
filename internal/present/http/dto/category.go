package dto

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
)

type CreateCategoryRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=255"`
	Slug    string `json:"slug" validate:"required,min=1,max=255"`
	ImageID *uint  `json:"image_id"`
}

func (r *CreateCategoryRequest) ToDomain() *domain.Category {
	return &domain.Category{
		Name: r.Name,
		Slug: r.Slug,
		// ImageID will be handled in service if needed, or domain struct updated to hold ImageID
	}
}

type UpdateCategoryRequest struct {
	ID      uint    `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Slug    *string `json:"slug,omitempty"`
	ImageID *uint   `json:"image_id,omitempty"`
}

type CategoryResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Image     *ImageResponse `json:"image,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

func NewCategoryResponse(category *domain.Category) *CategoryResponse {
	var imageResp *ImageResponse
	if category.Image != nil {
		imageResp = NewImageResponse(category.Image)
	}

	return &CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Image:     imageResp,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
