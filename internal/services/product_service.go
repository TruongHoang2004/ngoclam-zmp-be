package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type ProductService struct {
	productRepository *repositories.ProductRepository
	imageRepository   *repositories.ImageRepository
}

func NewProductService(productRepo *repositories.ProductRepository, imageRepo *repositories.ImageRepository) *ProductService {
	return &ProductService{
		productRepository: productRepo,
		imageRepository:   imageRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) *common.Error {

	varNameSet := make(map[string]struct{})
	for _, v := range product.Variants {
		name := v.Name
		if name == "" {
			continue // optional: skip empty names
		}
		if _, exists := varNameSet[name]; exists {
			return common.ErrConflict(ctx, "Product variant", "already exists")
		}
		varNameSet[name] = struct{}{}
	}

	domainProduct := product.ToDomain()

	existed, err := s.productRepository.IsExistProduct(ctx, domainProduct.Name)
	if err != nil {
		return err
	}

	if existed {
		return common.ErrConflict(ctx, "Product", "already exists")
	}

	if domainProduct.Variants == nil {
		dv := make([]domain.ProductVariant, 0, len(product.Variants))
		for _, v := range product.Variants {
			dv = append(dv, domain.ProductVariant{
				Name:  v.Name,
				Price: v.Price,
				Stock: v.Stock,
			})
		}
		domainProduct.Variants = &dv
	}

	return s.productRepository.CreateProduct(ctx, domainProduct)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*domain.Product, *common.Error) {
	product, err := s.productRepository.GetProductDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, page int, size int) ([]*domain.Product, int64, *common.Error) {
	offset := (page - 1) * size
	return s.productRepository.ListProducts(ctx, offset, size)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *dto.UpdateProductRequest) *common.Error {
	productdomain, err := s.productRepository.GetProductByID(ctx, product.ID)
	if err != nil {
		return common.ErrNotFound(ctx, "Product", "not found")
	}

	if product.Name != nil {
		productdomain.Name = *product.Name
	}
	if product.Description != nil {
		productdomain.Description = product.Description
	}
	if product.Price != nil {
		productdomain.Price = *product.Price
	}

	return s.productRepository.UpdateProduct(ctx, productdomain)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) *common.Error {
	productdomain, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepository.DeleteProduct(ctx, productdomain.ID)
}

func (s *ProductService) AddProductVariant(ctx context.Context, variant *dto.AddProductVariantRequest) *common.Error {
	product, err := s.productRepository.GetProductByID(ctx, variant.ProductID)
	if product == nil {
		return common.ErrNotFound(ctx, "Product", "not found")
	}
	if err != nil {
		return err
	}

	variantdomain := &domain.ProductVariant{
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Price:     variant.Price,
		Stock:     variant.Stock,
	}

	existed, err := s.productRepository.IsExistProductVariant(ctx, variant.ProductID, variantdomain.Name)
	if err != nil {
		return err
	}
	if existed {
		return common.ErrConflict(ctx, "Product variant", "already exists")
	}

	return s.productRepository.AddProductVariant(ctx, variantdomain)
}

func (s *ProductService) UpdateProductVariant(ctx context.Context, variant *dto.UpdateProductVariantRequest) *common.Error {

	variantdomain, err := s.productRepository.GetProductVariantByID(ctx, variant.ID)
	if err != nil {
		return err
	}

	if variant.Name != nil {
		variantdomain.Name = *variant.Name
	}
	if variant.Price != nil {
		variantdomain.Price = *variant.Price
	}
	if variant.Stock != nil {
		variantdomain.Stock = *variant.Stock
	}

	return s.productRepository.UpdateProductVariant(ctx, variantdomain)
}

func (s *ProductService) DeleteProductVariant(ctx context.Context, id uint) *common.Error {
	variantdomain, err := s.productRepository.GetProductVariantByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepository.DeleteProductVariant(ctx, variantdomain.ID)
}

func (s *ProductService) ListProductImages(ctx context.Context, productID uint) ([]*domain.ProductImage, *common.Error) {
	_, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return s.productRepository.ListProductImages(ctx, productID)
}

func (s *ProductService) AddProductImage(ctx context.Context, productID uint, req *dto.AttachProductImageRequest) (*domain.ProductImage, *common.Error) {
	_, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	_, err = s.imageRepository.GetImageByID(ctx, req.ImageID)
	if err != nil {
		return nil, err
	}

	order := 0
	if req.Order != nil {
		order = *req.Order
	}

	productImage := &domain.ProductImage{
		ProductID: productID,
		ImageID:   req.ImageID,
		Order:     order,
		IsMain:    req.IsMain,
	}

	return s.productRepository.AddProductImage(ctx, productImage)
}

func (s *ProductService) UpdateProductImage(ctx context.Context, productID uint, productImageID uint, req *dto.UpdateProductImageRequest) (*domain.ProductImage, *common.Error) {
	productImage, err := s.productRepository.GetProductImageByID(ctx, productImageID)
	if err != nil {
		return nil, err
	}

	if req.Order != nil {
		productImage.Order = *req.Order
	}

	if req.IsMain != nil {
		productImage.IsMain = *req.IsMain
	}

	return s.productRepository.UpdateProductImage(ctx, productImage)
}

func (s *ProductService) DeleteProductImage(ctx context.Context, productID uint, productImageID uint) *common.Error {
	productImage, err := s.productRepository.GetProductImageByID(ctx, productImageID)
	if err != nil {
		return err
	}
	if productImage == nil || productImage.ProductID != productID {
		return common.ErrNotFound(ctx, "Image", "not found")
	}

	return s.productRepository.DeleteProductImage(ctx, productImageID)
}
