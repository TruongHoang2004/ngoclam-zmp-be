package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
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

	existed, err := s.productRepository.IsExistProduct(ctx, product.Name)
	if err != nil {
		return err
	}

	if existed {
		return common.ErrConflict(ctx, "Product", "already exists")
	}

	newProduct := product.ToModel()
	newProduct.CategoryID = product.CategoryID

	if product.Variants != nil {
		var variants []model.ProductVariant
		for _, v := range product.Variants {
			variants = append(variants, model.ProductVariant{
				Name:  v.Name,
				Price: v.Price,
				Stock: v.Stock,
			})
		}
		newProduct.Variants = variants
	}

	return s.productRepository.CreateProduct(ctx, newProduct)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*model.Product, *common.Error) {
	product, err := s.productRepository.GetProductDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, page int, size int) ([]*model.Product, int64, *common.Error) {
	offset := (page - 1) * size
	return s.productRepository.ListProducts(ctx, offset, size)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *dto.UpdateProductRequest) *common.Error {
	productmodel, err := s.productRepository.GetProductByID(ctx, product.ID)
	if err != nil {
		return common.ErrNotFound(ctx, "Product", "not found")
	}

	if product.Name != nil {
		productmodel.Name = *product.Name
	}
	if product.Description != nil {
		productmodel.Description = product.Description
	}
	if product.Price != nil {
		productmodel.Price = *product.Price
	}

	return s.productRepository.UpdateProduct(ctx, productmodel)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) *common.Error {
	productmodel, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepository.DeleteProduct(ctx, productmodel.ID)
}

func (s *ProductService) AddProductVariant(ctx context.Context, req *dto.AddProductVariantRequest) *common.Error {
	product, err := s.productRepository.GetProductByID(ctx, req.ProductID)
	if product == nil {
		return common.ErrNotFound(ctx, "Product", "not found")
	}
	if err != nil {
		return err
	}

	variantmodel := &model.ProductVariant{
		ProductID: req.ProductID,
		Name:      req.Name,
		Price:     req.Price,
		Stock:     req.Stock,
	}

	existed, err := s.productRepository.IsExistProductVariant(ctx, req.ProductID, variantmodel.Name)
	if err != nil {
		return err
	}
	if existed {
		return common.ErrConflict(ctx, "Product variant", "already exists")
	}

	return s.productRepository.AddProductVariant(ctx, variantmodel)
}

func (s *ProductService) UpdateProductVariant(ctx context.Context, req *dto.UpdateProductVariantRequest) (*model.ProductVariant, *common.Error) {

	variantmodel, err := s.productRepository.GetProductVariantByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		variantmodel.Name = *req.Name
	}
	if req.Price != nil {
		variantmodel.Price = *req.Price
	}
	if req.Stock != nil {
		variantmodel.Stock = *req.Stock
	}

	err = s.productRepository.UpdateProductVariant(ctx, variantmodel)
	if err != nil {
		return nil, err
	}

	return variantmodel, nil
}

func (s *ProductService) DeleteProductVariant(ctx context.Context, id uint) *common.Error {
	variantmodel, err := s.productRepository.GetProductVariantByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepository.DeleteProductVariant(ctx, variantmodel.ID)
}

func (s *ProductService) ListProductImages(ctx context.Context, productID uint) ([]*model.ProductImage, *common.Error) {
	_, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return s.productRepository.ListProductImages(ctx, productID)
}

func (s *ProductService) AddProductImage(ctx context.Context, productID uint, req *dto.AttachProductImageRequest) (*model.ProductImage, *common.Error) {
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

	productImage := &model.ProductImage{
		ProductID: productID,
		ImageID:   req.ImageID,
		Order:     order,
		IsMain:    req.IsMain,
	}

	return s.productRepository.AddProductImage(ctx, productImage)
}

func (s *ProductService) UpdateProductImage(ctx context.Context, productID uint, productImageID uint, req *dto.UpdateProductImageRequest) (*model.ProductImage, *common.Error) {
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
