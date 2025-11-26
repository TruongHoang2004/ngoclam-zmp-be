package domain

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Image     *Image    `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategoryDomain(category *model.Category) *Category {
	return &Category{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (c *Category) ToModel() *model.Category {
	imageID := uint(0)
	if c.Image != nil {
		imageID = c.Image.ID
	}

	return &model.Category{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		ImageID:   &imageID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c *Category) AddImage(image *model.Image) {
	c.Image = &Image{
		ID:        image.ID,
		Name:      image.Name,
		URL:       image.URL,
		Hash:      image.Hash,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}
