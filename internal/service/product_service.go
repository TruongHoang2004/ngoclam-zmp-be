package service

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) error {
	return s.productRepository.CreateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	return s.productRepository.GetProductByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) error {
	return s.productRepository.UpdateProduct(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepository.DeleteProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, offset int, limit int) ([]*model.Product, error) {
	return s.productRepository.ListProducts(ctx, offset, limit)
}
