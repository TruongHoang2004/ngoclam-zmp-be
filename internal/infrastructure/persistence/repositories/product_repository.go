package repositories

import (
	"context"
	"errors"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	*baseRepository
}

func NewProductRepository(base *baseRepository) *ProductRepository {
	return &ProductRepository{baseRepository: base}
}

// CreateProduct creates a new product (variants should be added separately)
func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) *common.Error {
	return r.returnError(ctx, r.db.WithContext(ctx).Create(product).Error)
}

// IsExistProduct checks if a product with the same name already exists
func (r *ProductRepository) IsExistProduct(ctx context.Context, name string) (bool, *common.Error) {
	var count int64

	conds := []clause.Expression{
		clause.Eq{Column: "name", Value: name},
	}

	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Clauses(conds...).
		Count(&count).
		Error

	if err != nil {
		return false, common.ErrSystemError(ctx, err.Error())
	}

	return count > 0, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint) (*model.Product, *common.Error) {
	product, err := r.GetProductDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetProductSummaryByID(ctx context.Context, id uint) (*model.Product, *common.Error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).
		Preload("ProductImages", "is_main = ?", true).
		Preload("ProductImages.Image").
		First(&prod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, r.returnError(ctx, err)
	}

	return &prod, nil
}

func (r *ProductRepository) GetProductDetailByID(ctx context.Context, id uint) (*model.Product, *common.Error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).
		Preload("Variants").
		Preload("ProductImages.Image").
		First(&prod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, r.returnError(ctx, err)
	}

	return &prod, nil
}

// ListProducts with pagination, always includes variants
func (r *ProductRepository) ListProducts(ctx context.Context, offset, limit int) ([]*model.Product, int64, *common.Error) {
	// Count
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	var products []*model.Product
	if err := r.db.WithContext(ctx).
		Preload("ProductImages", "is_main = ?", true).
		Preload("ProductImages.Image").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return products, total, nil
}

func (r *ProductRepository) GetProductsByCategoryID(ctx context.Context, categoryID uint, offset, limit int) ([]*model.Product, int64, *common.Error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	var products []*model.Product
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("ProductImages", "is_main = ?", true).
		Preload("ProductImages.Image").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return products, total, nil
}

// UpdateProduct updates only product fields (not variants)
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) *common.Error {
	m := &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
	}
	err := r.db.WithContext(ctx).Model(m).Updates(m).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uint) *common.Error {
	err := r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

// === Product Variant Methods ===

func (r *ProductRepository) GetProductVariantByID(ctx context.Context, id uint) (*model.ProductVariant, *common.Error) {
	var v model.ProductVariant
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product variant", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}
	return &v, nil
}

func (r *ProductRepository) AddProductVariant(ctx context.Context, variant *model.ProductVariant) *common.Error {
	err := r.db.WithContext(ctx).Create(variant).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductVariant(ctx context.Context, variant *model.ProductVariant) *common.Error {
	m := &model.ProductVariant{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Stock:     variant.Stock,
		Price:     variant.Price,
	}
	err := r.db.WithContext(ctx).Model(m).Updates(m).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) DeleteProductVariant(ctx context.Context, id uint) *common.Error {
	err := r.db.WithContext(ctx).Delete(&model.ProductVariant{}, id).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) IsExistProductVariant(ctx context.Context, productID uint, name string) (bool, *common.Error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ProductVariant{}).
		Where("product_id = ? AND name = ?", productID, name).
		Count(&count).Error
	if err != nil {
		return false, r.returnError(ctx, err)
	}
	return count > 0, nil
}

// === Product Image Methods ===

func (r *ProductRepository) ListProductImages(ctx context.Context, productID uint) ([]*model.ProductImage, *common.Error) {
	images, err := r.loadProductImagesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ProductRepository) GetProductImageByID(ctx context.Context, id uint) (*model.ProductImage, *common.Error) {
	var record model.ProductImage
	if err := r.db.WithContext(ctx).
		First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Image", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}

	var img []model.Image
	if err := r.db.WithContext(ctx).
		First(&img, record.ImageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Image", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}

	return &record, nil
}

func (r *ProductRepository) AddProductImage(ctx context.Context, img *model.ProductImage) (*model.ProductImage, *common.Error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if img.IsMain {
			if err := tx.Model(&model.ProductImage{}).
				Where("product_id = ?", img.ProductID).
				Update("is_main", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(img).Error; err != nil {
			return err
		}

		return tx.Preload("Image").
			Preload("Variant").
			First(img, img.ID).Error
	})
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	return img, nil
}

func (r *ProductRepository) UpdateProductImage(ctx context.Context, img *model.ProductImage) (*model.ProductImage, *common.Error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if img.IsMain {
			if err := tx.Model(&model.ProductImage{}).
				Where("product_id = ?", img.ProductID).
				Update("is_main", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&model.ProductImage{}).
			Where("id = ?", img.ID).
			Updates(map[string]interface{}{
				"order":   img.Order,
				"is_main": img.IsMain,
			}).Error; err != nil {
			return err
		}

		return tx.Preload("Image").
			Preload("Variant").
			First(img, img.ID).Error
	})
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	return img, nil
}

func (r *ProductRepository) DeleteProductImage(ctx context.Context, id uint) *common.Error {
	err := r.db.WithContext(ctx).Delete(&model.ProductImage{}, id).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) loadProductImagesByProductID(ctx context.Context, productID uint) ([]*model.ProductImage, *common.Error) {
	var images []*model.ProductImage
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("is_main DESC, \"order\" ASC, id ASC").
		Find(&images).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}

	return images, nil
}
