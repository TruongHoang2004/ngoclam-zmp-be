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

func (r *ProductRepository) IsExistProduct(ctx context.Context, product *model.Product) (bool, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Where("name = ?", product.Name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAllProductVariantsByProductID(ctx context.Context, productID uint) ([]model.ProductVariant, error) {
	var variants []model.ProductVariant
	if err := r.db.Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *ProductRepository) IsExistProductvariants(ctx context.Context, productID uint, variant *model.ProductVariant) (bool, error) {
	var count int64
	err := r.db.Model(&model.ProductVariant{}).
		Where("product_id = ? AND name = ?", productID, variant.Name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) ListProducts(ctx context.Context, offset int, limit int) ([]*model.Product, int64, error) {
	query := r.db.Model(&model.Product{})
	var products []*model.Product
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	query.Count(&total)
	return products, total, nil
}

func (r *ProductRepository) AddProductVariant(ctx context.Context, variant *model.ProductVariant) error {
	return r.db.Create(variant).Error
}
