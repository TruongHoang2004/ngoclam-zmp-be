package service

import (
	"context"
	"fmt"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) error {

	varNameSet := make(map[string]struct{})
	for _, v := range product.Variants {
		name := v.Name
		if name == "" {
			continue // optional: skip empty names
		}
		if _, exists := varNameSet[name]; exists {
			return common.BadRequest(fmt.Sprintf("Duplicate variant name: %s", name))
		}
		varNameSet[name] = struct{}{}
	}

	domainProduct := product.ToDomain()

	existed, err := s.productRepository.IsExistProduct(ctx, domainProduct.Name)
	if err != nil {
		return err
	}

	if existed {
		return common.Conflict("Product already exists")
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

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	product, err := s.productRepository.GetProductDetailByID(ctx, id, true)
	if err != nil {
		return nil, common.NotFound("Product not found")
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, page int, size int) ([]*domain.Product, int64, error) {
	offset := (page - 1) * size
	return s.productRepository.ListProducts(ctx, offset, size)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *dto.UpdateProductRequest) error {
	productdomain, err := s.productRepository.GetProductDetailByID(ctx, product.ID, false)
	if err != nil {
		return common.NotFound("Product not found")
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

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	productdomain, err := s.productRepository.GetProductDetailByID(ctx, id, false)
	if err != nil {
		return common.NotFound("Product not found")
	}

	return s.productRepository.DeleteProduct(ctx, productdomain.ID)
}

func (s *ProductService) AddProductVariant(ctx context.Context, variant *dto.AddProductVariantRequest) error {
	product, err := s.productRepository.GetProductDetailByID(ctx, variant.ProductID, false)
	if product == nil {
		return common.NotFound("Product not found")
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
		return common.Conflict("Product variant already exists")
	}

	return s.productRepository.AddProductVariant(ctx, variantdomain)
}

func (s *ProductService) UpdateProductVariant(ctx context.Context, variant *dto.UpdateProductVariantRequest) error {

	variantdomain := s.productRepository.GetProductVariantByID(ctx, variant.ID)
	if variantdomain == nil {
		return common.NotFound("Product variant not found")
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

func (s *ProductService) DeleteProductVariant(ctx context.Context, id uint) error {
	variantdomain := s.productRepository.GetProductVariantByID(ctx, id)
	if variantdomain == nil {
		return common.NotFound("Product variant not found")
	}

	return s.productRepository.DeleteProductVariant(ctx, variantdomain.ID)
}
