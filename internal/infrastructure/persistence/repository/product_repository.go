package repository

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) ListProducts(ctx context.Context, offset int, limit int) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
