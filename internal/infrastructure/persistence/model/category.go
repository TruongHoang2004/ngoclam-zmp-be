package model

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
	"gorm.io/gorm"
)

type CategoryModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`

	ImageRelated *ImageRelatedModel `gorm:"polymorphic:Entity;polymorphicValue:category"`
	Products     *[]Product         `gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

func MapCategoryToModel(category *entity.Category) *CategoryModel {

	return &CategoryModel{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ImageRelated: &ImageRelatedModel{
			ImageID:    category.Image.ID,
			EntityID:   category.ID,
			EntityType: EntityTypeCategory,
			Order:      0,
		},
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (c *CategoryModel) ToDomain() *entity.Category {
	var image *entity.Image = nil
	if c.ImageRelated != nil {
		image = c.ImageRelated.Image.ToDomain()
	}
	return &entity.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Image:       image,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
