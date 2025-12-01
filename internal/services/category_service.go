package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type CategoryService struct {
	*baseService
	categoryRepository *repositories.CategoryRepository
	productRepository  *repositories.ProductRepository
}

func NewCategoryService(
	categoryRepository *repositories.CategoryRepository,
	productRepository *repositories.ProductRepository,
) *CategoryService {
	return &CategoryService{
		baseService:        NewBaseService(),
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) *common.Error {
	// Implementation for creating a category

	if s.categoryRepository.IsExist(ctx, req.Name, req.Slug) {
		return common.ErrConflict(ctx, "Category", "Category already exists")
	}

	category := req.ToModel()
	if req.ImageID != nil {
		category.ImageID = req.ImageID
	}
	return s.categoryRepository.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*model.Category, *common.Error) {
	// Implementation for retrieving a category by ID

	return s.categoryRepository.GetCategoryByID(ctx, id)
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*model.Category, *common.Error) {
	return s.categoryRepository.ListCategories(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*model.Category, *common.Error) {
	// Implementation for updating a category

	category, err := s.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.ImageID != nil {
		category.ImageID = req.ImageID
	}

	err = s.categoryRepository.UpdateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) *common.Error {
	return s.categoryRepository.DeleteCategory(ctx, id)
}

func (s *CategoryService) GetProductsByCategory(ctx context.Context, categoryID uint, page, size int) (*dto.PaginationResponse[dto.ProductResponse], *common.Error) {
	products, total, err := s.productRepository.GetProductsByCategoryID(ctx, categoryID, (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *dto.NewProductResponse(p))
	}

	response := dto.NewPaginationResponse(productResponses, total, dto.PaginationRequest{Page: page, Size: size})
	return &response, nil
}
