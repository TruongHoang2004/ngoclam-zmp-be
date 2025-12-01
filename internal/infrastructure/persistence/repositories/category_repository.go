package repositories

import (
	"context"
	"errors"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	// Repository methods would be defined here
	*baseRepository
}

func NewCategoryRepository(base *baseRepository) *CategoryRepository {
	return &CategoryRepository{baseRepository: base}
}

func (c *CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) *common.Error {
	// Implementation for creating a category

	return c.returnError(ctx, c.db.WithContext(ctx).Create(category).Error)
}

func (c *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*model.Category, *common.Error) {
	// Implementation for retrieving a category by ID

	var model model.Category
	err := c.db.WithContext(ctx).Preload("Image").First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, c.returnError(ctx, common.ErrNotFound(ctx, "Category", "Notfound"))
		}
	}
	return &model, nil
}

func (c *CategoryRepository) IsExist(ctx context.Context, name, slug string) bool {
	nameCount := int64(0)
	slugCount := int64(0)
	c.db.WithContext(ctx).Model(&model.Category{}).Where("name = ?", name).Count(&nameCount)
	c.db.WithContext(ctx).Model(&model.Category{}).Where("slug = ?", slug).Count(&slugCount)
	return nameCount > 0 || slugCount > 0
}

func (c *CategoryRepository) GetCategoryDetail(ctx context.Context, id uint) (*model.Category, *common.Error) {
	var categoryModel model.Category
	err := c.db.WithContext(ctx).Preload("Image").First(&categoryModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Category", "Not found")
		}
	}
	return &categoryModel, nil

}

func (c *CategoryRepository) ListCategories(ctx context.Context) ([]*model.Category, *common.Error) {
	var models []*model.Category
	err := c.db.WithContext(ctx).Preload("Image").Find(&models).Error
	if err != nil {
		return nil, c.returnError(ctx, err)
	}

	return models, nil
}

func (c *CategoryRepository) UpdateCategory(ctx context.Context, category *model.Category) *common.Error {
	// Implementation for updating a category

	return c.returnError(ctx, c.db.WithContext(ctx).Save(category).Error)
}

func (c *CategoryRepository) DeleteCategory(ctx context.Context, id uint) *common.Error {
	// Implementation for deleting a category

	productNumber := int64(0)
	err := c.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", id).Count(&productNumber).Error

	if err != nil {
		return c.returnError(ctx, err)
	}

	if productNumber > 0 {
		return c.returnError(ctx, common.ErrConflict(ctx, "Category", "Cannot delete category with associated products"))
	}

	return c.returnError(ctx, c.db.WithContext(ctx).Delete(&model.Category{}, id).Error)
}
