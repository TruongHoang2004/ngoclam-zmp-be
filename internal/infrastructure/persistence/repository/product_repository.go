package repository

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct creates a new product (variants should be added separately)
func (r *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	m := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// IsExistProduct checks if a product with the same name already exists
func (r *ProductRepository) IsExistProduct(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("name = ?", name).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetProductDetailByID fetches product + manually loads variants if requested
func (r *ProductRepository) GetProductDetailByID(ctx context.Context, id uint, withVariants bool) (*domain.Product, error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		return nil, err
	}

	domainProd := domain.NewProductFromModel(&prod)

	if withVariants {
		variants, err := r.loadVariantsByProductID(ctx, id)
		if err != nil {
			return nil, err
		}
		domainProd.Variants = &variants
	}

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

	// Load all variants in one query
	variantMap := make(map[uint][]domain.ProductVariant)
	if len(ids) > 0 {
		var variants []model.ProductVariant
		if err := r.db.WithContext(ctx).
			Where("product_id IN ?", ids).
			Find(&variants).Error; err != nil {
			return nil, 0, err
		}

		for _, v := range variants {
			dv := domain.ProductVariant{
				ID:        v.ID,
				ProductID: v.ProductID,
				Name:      v.Name,
				Stock:     v.Stock,
				Price:     v.Price,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			}
			variantMap[v.ProductID] = append(variantMap[v.ProductID], dv)
		}
	}

	// Build final domain products
	result := make([]*domain.Product, len(products))
	for i, p := range products {
		domainProd := domain.NewProductFromModel(p)
		if vars, ok := variantMap[p.ID]; ok && len(vars) > 0 {
			domainProd.Variants = &vars
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
