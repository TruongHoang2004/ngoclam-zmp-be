package repository

import (
	"context"
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) entity.CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (c *CategoryRepositoryImpl) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	categoryModel := model.MapCategoryToModel(category)

	log.Printf("Category Model: %v", categoryModel.ImageRelated)

	if err := c.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(categoryModel).Error; err != nil {
		return nil, err
	}

	// Load the associated ImageRelated after creation
	if err := c.db.WithContext(ctx).Preload("ImageRelated").Preload("ImageRelated.Image").First(categoryModel, categoryModel.ID).Error; err != nil {
		return nil, err
	}

	// Convert back to domain entity

	return categoryModel.ToDomain(), nil
}

// FindByID returns category by ID with optional image
func (c *CategoryRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Category, error) {
	var categoryModel model.CategoryModel
	if err := c.db.WithContext(ctx).Preload("ImageRelated").Preload("ImageRelated.Image").First(&categoryModel, id).Error; err != nil {
		return nil, err
	}

	category := categoryModel.ToDomain()
	return category, nil
}

// FindAll returns all categories with optional images
func (c *CategoryRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Category, error) {
	var categoryModels []model.CategoryModel
	if err := c.db.WithContext(ctx).Preload("ImageRelated").Preload("ImageRelated.Image").Find(&categoryModels).Error; err != nil {
		return nil, err
	}

	var categories []*entity.Category
	for _, categoryModel := range categoryModels {

		categories = append(categories, categoryModel.ToDomain())
	}
	return categories, nil
}

// Update modifies an existing category
func (c *CategoryRepositoryImpl) Update(ctx context.Context, category *entity.Category) error {
	categoryModel := model.MapCategoryToModel(category)

	return c.db.WithContext(ctx).Save(categoryModel).Error
}

// Delete removes a category by ID
func (c *CategoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).Delete(&model.CategoryModel{}, id).Error
}
