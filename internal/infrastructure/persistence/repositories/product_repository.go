package repositories

import (
	"context"
	"errors"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
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
func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) *common.Error {
	m := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	return r.returnError(ctx, r.db.WithContext(ctx).Create(m).Error)
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

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	return r.GetProductDetailByID(ctx, id)
}

func (r *ProductRepository) GetProductSummaryByID(ctx context.Context, id uint) (*domain.Product, error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		return nil, err
	}

	domainProd := domain.NewProductFromModel(&prod)

	var img model.ProductImage
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_main = ?", id, true).
		Preload("Image").
		Preload("Variant").
		First(&img).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainProd, nil
		}
		return nil, err
	}

	domainProd.SetImages([]*model.ProductImage{&img})

	return domainProd, nil
}

// GetProductDetailByID fetches product + manually loads variants if requested
func (r *ProductRepository) GetProductDetailByID(ctx context.Context, id uint) (*domain.Product, error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		return nil, err
	}

	domainProd := domain.NewProductFromModel(&prod)

	variants, err := r.loadVariantsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	domainProd.Variants = &variants

	images, err := r.loadProductImagesByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	domainProd.SetImages(images)

	return domainProd, nil
}

// ListProducts with pagination, always includes variants
func (r *ProductRepository) ListProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []*model.Product
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	// Collect all product IDs
	ids := make([]uint, 0, len(products))
	for _, p := range products {
		if p != nil {
			ids = append(ids, p.ID)
		}
	}

	imageMap := make(map[uint][]domain.ProductImage)
	if len(ids) > 0 {
		var productImages []*model.ProductImage
		if err := r.db.WithContext(ctx).
			Where("product_id IN ? AND is_main = ?", ids, true).
			Preload("Image").
			Preload("Variant").
			Find(&productImages).Error; err != nil {
			return nil, 0, err
		}

		for _, img := range productImages {
			if img == nil {
				continue
			}
			// keep only one main image per product
			if _, exists := imageMap[img.ProductID]; exists {
				continue
			}
			if d := domain.NewProductImageFromModel(img); d != nil {
				imageMap[img.ProductID] = []domain.ProductImage{*d}
			}
		}
	}

	// Build final domain products
	result := make([]*domain.Product, len(products))
	for i, p := range products {
		domainProd := domain.NewProductFromModel(p)
		if imgs, ok := imageMap[p.ID]; ok && len(imgs) > 0 {
			domainProd.Images = &imgs
		}
		result[i] = domainProd
	}

	return result, total, nil
}

// UpdateProduct updates only product fields (not variants)
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	m := &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// === Product Variant Methods ===

func (r *ProductRepository) GetProductVariantByID(ctx context.Context, id uint) *domain.ProductVariant {
	var v model.ProductVariant
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil
	}
	return &domain.ProductVariant{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		Stock:     v.Stock,
		Price:     v.Price,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func (r *ProductRepository) AddProductVariant(ctx context.Context, variant *domain.ProductVariant) error {
	m := &model.ProductVariant{
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Stock:     variant.Stock,
		Price:     variant.Price,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ProductRepository) UpdateProductVariant(ctx context.Context, variant *domain.ProductVariant) error {
	m := &model.ProductVariant{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Stock:     variant.Stock,
		Price:     variant.Price,
	}
	return r.db.WithContext(ctx).Model(m).Updates(m).Error
}

func (r *ProductRepository) DeleteProductVariant(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ProductVariant{}, id).Error
}

func (r *ProductRepository) IsExistProductVariant(ctx context.Context, productID uint, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ProductVariant{}).
		Where("product_id = ? AND name = ?", productID, name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Helper: load variants for a single product
func (r *ProductRepository) loadVariantsByProductID(ctx context.Context, productID uint) ([]domain.ProductVariant, error) {
	var variants []model.ProductVariant
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, err
	}

	result := make([]domain.ProductVariant, len(variants))
	for i, v := range variants {
		result[i] = domain.ProductVariant{
			ID:        v.ID,
			ProductID: v.ProductID,
			Name:      v.Name,
			Stock:     v.Stock,
			Price:     v.Price,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}
	return result, nil
}

// === Product Image Methods ===

func (r *ProductRepository) ListProductImages(ctx context.Context, productID uint) ([]*domain.ProductImage, error) {
	images, err := r.loadProductImagesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.ProductImage, 0, len(images))
	for _, img := range images {
		if img == nil {
			continue
		}
		if d := domain.NewProductImageFromModel(img); d != nil {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *ProductRepository) GetProductImageByID(ctx context.Context, id uint) (*domain.ProductImage, error) {
	var record model.ProductImage
	if err := r.db.WithContext(ctx).
		Preload("Image").
		Preload("Variant").
		First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewProductImageFromModel(&record), nil
}

func (r *ProductRepository) AddProductImage(ctx context.Context, img *domain.ProductImage) (*domain.ProductImage, error) {
	record := img.ToModel()
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if record.IsMain {
			if err := tx.Model(&model.ProductImage{}).
				Where("product_id = ?", record.ProductID).
				Update("is_main", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return tx.Preload("Image").
			Preload("Variant").
			First(record, record.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return domain.NewProductImageFromModel(record), nil
}

func (r *ProductRepository) UpdateProductImage(ctx context.Context, img *domain.ProductImage) (*domain.ProductImage, error) {
	record := img.ToModel()
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if record.IsMain {
			if err := tx.Model(&model.ProductImage{}).
				Where("product_id = ?", record.ProductID).
				Update("is_main", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&model.ProductImage{}).
			Where("id = ?", record.ID).
			Updates(map[string]interface{}{
				"order":   record.Order,
				"is_main": record.IsMain,
			}).Error; err != nil {
			return err
		}

		return tx.Preload("Image").
			Preload("Variant").
			First(record, record.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return domain.NewProductImageFromModel(record), nil
}

func (r *ProductRepository) DeleteProductImage(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ProductImage{}, id).Error
}

func (r *ProductRepository) loadProductImagesByProductID(ctx context.Context, productID uint) ([]*model.ProductImage, error) {
	var images []*model.ProductImage
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Preload("Image").
		Preload("Variant").
		Order("is_main DESC, \"order\" ASC, id ASC").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
