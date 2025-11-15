package service

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
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
	modelProduct := &model.Product{
		Name:        product.Name,
		Description: &product.Description,
		Price:       product.Price,
	}

	existed, err := s.productRepository.IsExistProduct(ctx, modelProduct)
	if err != nil {
		return err
	}
	if existed {
		return common.Conflict("Product already exists")
	}

	return s.productRepository.CreateProduct(ctx, modelProduct)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*model.Product, []model.ProductVariant, error) {
	product, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, nil, common.NotFound("Product not found")
	}

	variants, err := s.productRepository.GetAllProductVariantsByProductID(ctx, product.ID)
	if err != nil {
		return nil, nil, err
	}

	return product, variants, nil
}

func (s *ProductService) ListProducts(ctx context.Context, page int, size int) ([]*model.Product, int64, error) {
	offset := (page - 1) * size
	return s.productRepository.ListProducts(ctx, offset, size)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *dto.UpdateProductRequest) error {
	productModel, err := s.productRepository.GetProductByID(ctx, product.ID)
	if err != nil {
		return common.NotFound("Product not found")
	}

	if product.Name != nil {
		productModel.Name = *product.Name
	}
	if product.Description != nil {
		productModel.Description = product.Description
	}
	if product.Price != nil {
		productModel.Price = *product.Price
	}

	return s.productRepository.UpdateProduct(ctx, productModel)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	productModel, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return common.NotFound("Product not found")
	}

	return s.productRepository.DeleteProduct(ctx, productModel.ID)
}

func (s *ProductService) AddProductVariant(ctx context.Context, variant *dto.AddProductVariantRequest) error {
	product, err := s.productRepository.GetProductByID(ctx, variant.ProductID)
	if product == nil {
		return common.NotFound("Product not found")
	}
	if err != nil {
		return err
	}

	variantModel := &model.ProductVariant{
		ProductID: variant.ProductID,
		Name:      variant.Name,
		Price:     variant.Price,
		Stock:     variant.Stock,
	}

	existed, err := s.productRepository.IsExistProductvariants(ctx, variant.ProductID, variantModel)
	if err != nil {
		return err
	}
	if existed {
		return common.Conflict("Product variant already exists")
	}

	return s.productRepository.AddProductVariant(ctx, variantModel)
}
