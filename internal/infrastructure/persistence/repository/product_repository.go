package repository

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) entity.ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

// Add methods for ProductRepository as needed
func (p *ProductRepositoryImpl) Create(ctx context.Context, product entity.Product) (*entity.Product, error) {
	productModel := model.MapProductToModel(&product)

	tx := p.db.WithContext(ctx).Begin()

	for i := range productModel.Images {
		productModel.Images[i].EntityType = model.EntityTypeProduct
	}

	if err := tx.Create(productModel).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return productModel.ToDomain(), nil
}

func (p *ProductRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	var productModel model.Product
	if err := p.db.WithContext(ctx).Preload("Category.ImageRelated").Preload("Images").Preload("Variants").Preload("Images.Image").Preload("Variants.Image").First(&productModel, id).Error; err != nil {
		return nil, err
	}
	var imageModels []model.ImageModel
	if productModel.Images != nil {
		imageIDs := make([]uint, len(productModel.Images))
		for i, img := range productModel.Images {
			imageIDs[i] = img.ID
		}
		if err := p.db.Where("id IN ?", imageIDs).Find(&imageModels).Error; err != nil {
			return nil, err
		}
	}
	// We're not loading the category image as it's not needed
	// Just keep the product model with its existing preloaded relationships

	return productModel.ToDomain(), nil
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Product, error) {
	var productModels []model.Product
	if err := p.db.WithContext(ctx).Preload("Category.ImageRelated").Preload("Images.Image").Preload("Variants").Preload("Variants.Image").Find(&productModels).Error; err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(productModels))
	for i, pm := range productModels {
		var imageModels []model.ImageModel
		if pm.Images != nil {
			imageIDs := make([]uint, len(pm.Images))
			for j, img := range pm.Images {
				imageIDs[j] = img.ID
			}
			if err := p.db.Where("id IN ?", imageIDs).Find(&imageModels).Error; err != nil {
				return nil, err
			}
		}
		products[i] = pm.ToDomain()
	}
	return products, nil
}

func (p *ProductRepositoryImpl) Update(ctx context.Context, product entity.Product) (*entity.Product, error) {
	productModel := model.MapProductToModel(&product)
	if err := p.db.WithContext(ctx).Save(productModel).Error; err != nil {
		return nil, err
	}
	return productModel.ToDomain(), nil
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, id uint) error {
	if err := p.db.WithContext(ctx).Delete(&model.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID uint) ([]*entity.Product, error) {
	var productModels []model.Product
	if err := p.db.WithContext(ctx).Preload("Images").Preload("Variants").Where("category_id = ?", categoryID).Find(&productModels).Error; err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(productModels))
	for i, pm := range productModels {
		var imageModels []model.ImageModel
		if pm.Images != nil {
			imageIDs := make([]uint, len(pm.Images))
			for j, img := range pm.Images {
				imageIDs[j] = img.ID
			}
			if err := p.db.Where("id IN ?", imageIDs).Find(&imageModels).Error; err != nil {
				return nil, err
			}
		}
		products[i] = pm.ToDomain()
	}
	return products, nil
}
