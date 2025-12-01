package repositories

import (
	"context"
	"errors"
	"fmt"

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

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint) (*domain.Product, *common.Error) {
	product, err := r.GetProductDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetProductSummaryByID(ctx context.Context, id uint) (*domain.Product, *common.Error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, r.returnError(ctx, err)
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
		return nil, r.returnError(ctx, err)
	}

	domainProd.SetImages([]*model.ProductImage{&img})

	return domainProd, nil
}

// GetProductDetailByID fetches product + manually loads variants if requested
func (r *ProductRepository) GetProductDetailByID(ctx context.Context, id uint) (*domain.Product, *common.Error) {
	var prod model.Product
	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}
		return nil, r.returnError(ctx, err)
	}

	domainProd := domain.NewProductFromModel(&prod)

	variants, err := r.loadVariantsByProductID(ctx, id)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	domainProd.Variants = variants

	images, err := r.loadProductImagesByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	domainProd.SetImages(images)

	return domainProd, nil
}

// ListProducts with pagination, always includes variants
func (r *ProductRepository) ListProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int64, *common.Error) {
	// Count
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	// Struct chứa tuple kết quả JOIN
	type Row struct {
		Product      model.Product      `gorm:"embedded;embeddedPrefix:product_"` // tránh clash field
		ProductImage model.ProductImage `gorm:"embedded;embeddedPrefix:pi_"`
		Image        model.Image        `gorm:"embedded;embeddedPrefix:img_"`
	}

	var rows []Row

	err := r.db.WithContext(ctx).
		Table("products AS p").
		Select(
			"p.id AS product_id",
			"p.name AS product_name",
			"p.category_id AS product_category_id",
			"p.description AS product_description",
			"p.price AS product_price",
			"p.created_at AS product_created_at",
			"p.updated_at AS product_updated_at",

			"pi.id AS pi_id",
			"pi.product_id AS pi_product_id",
			"pi.image_id AS pi_image_id",
			"pi.is_main AS pi_is_main",
			"pi.created_at AS pi_created_at",
			"pi.updated_at AS pi_updated_at",

			"img.id AS img_id",
			"img.url AS img_url",
			"img.created_at AS img_created_at",
			"img.updated_at AS img_updated_at",
		).
		Joins("LEFT JOIN product_images AS pi ON pi.product_id = p.id AND pi.is_main = ?", true).
		Joins("LEFT JOIN images AS img ON img.id = pi.image_id").Offset(offset).
		Limit(limit).
		Scan(&rows).Error

	if err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	// Convert sang domain
	result := make([]*domain.Product, 0)

	for _, r := range rows {
		prod := domain.NewProductFromModel(&r.Product)
		fmt.Println(r.Product.Name)

		// Nếu có ảnh
		if r.ProductImage.ID != 0 {
			prod.SetImages([]*model.ProductImage{&r.ProductImage})
		}

		result = append(result, prod)
	}

	return result, total, nil
}

func (r *ProductRepository) GetProductsByCategoryID(ctx context.Context, categoryID uint, offset, limit int) ([]*domain.Product, int64, *common.Error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	type Row struct {
		Product      model.Product      `gorm:"embedded;embeddedPrefix:product_"`
		ProductImage model.ProductImage `gorm:"embedded;embeddedPrefix:pi_"`
		Image        model.Image        `gorm:"embedded;embeddedPrefix:img_"`
	}

	var rows []Row

	err := r.db.WithContext(ctx).
		Table("products AS p").
		Select(
			"p.id AS product_id",
			"p.name AS product_name",
			"p.category_id AS product_category_id",
			"p.description AS product_description",
			"p.price AS product_price",
			"p.created_at AS product_created_at",
			"p.updated_at AS product_updated_at",

			"pi.id AS pi_id",
			"pi.product_id AS pi_product_id",
			"pi.image_id AS pi_image_id",
			"pi.is_main AS pi_is_main",
			"pi.created_at AS pi_created_at",
			"pi.updated_at AS pi_updated_at",

			"img.id AS img_id",
			"img.url AS img_url",
			"img.created_at AS img_created_at",
			"img.updated_at AS img_updated_at",
		).
		Joins("LEFT JOIN product_images AS pi ON pi.product_id = p.id AND pi.is_main = ?", true).
		Joins("LEFT JOIN images AS img ON img.id = pi.image_id").
		Where("p.category_id = ?", categoryID).
		Offset(offset).
		Limit(limit).
		Scan(&rows).Error

	if err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	result := make([]*domain.Product, 0)

	for _, r := range rows {
		prod := domain.NewProductFromModel(&r.Product)

		if r.ProductImage.ID != 0 {
			prod.SetImages([]*model.ProductImage{&r.ProductImage})
		}

		result = append(result, prod)
	}

	return result, total, nil
}

// UpdateProduct updates only product fields (not variants)
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) *common.Error {
	m := &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
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

func (r *ProductRepository) GetProductVariantByID(ctx context.Context, id uint) (*domain.ProductVariant, *common.Error) {
	var v model.ProductVariant
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Product variant", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}
	return &domain.ProductVariant{
		ID:        v.ID,
		ProductID: v.ProductID,
		Name:      v.Name,
		Stock:     v.Stock,
		Price:     v.Price,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}, nil
}

func (r *ProductRepository) AddProductVariant(ctx context.Context, variant *domain.ProductVariant) *common.Error {
	m := &model.ProductVariant{
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Stock:     variant.Stock,
		Price:     variant.Price,
	}
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductVariant(ctx context.Context, variant *domain.ProductVariant) *common.Error {
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

// Helper: load variants for a single product
func (r *ProductRepository) loadVariantsByProductID(ctx context.Context, productID uint) (*[]domain.ProductVariant, *common.Error) {
	var variants []model.ProductVariant
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, r.returnError(ctx, err)
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
	return &result, nil
}

// === Product Image Methods ===

func (r *ProductRepository) ListProductImages(ctx context.Context, productID uint) ([]*domain.ProductImage, *common.Error) {
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

func (r *ProductRepository) GetProductImageByID(ctx context.Context, id uint) (*domain.ProductImage, *common.Error) {
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

	return domain.NewProductImageFromModel(&record), nil
}

func (r *ProductRepository) AddProductImage(ctx context.Context, img *domain.ProductImage) (*domain.ProductImage, *common.Error) {
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
		return nil, r.returnError(ctx, err)
	}
	return domain.NewProductImageFromModel(record), nil
}

func (r *ProductRepository) UpdateProductImage(ctx context.Context, img *domain.ProductImage) (*domain.ProductImage, *common.Error) {
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
		return nil, r.returnError(ctx, err)
	}
	return domain.NewProductImageFromModel(record), nil
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
