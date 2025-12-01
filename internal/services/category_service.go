package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
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

func (s *CategoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, *common.Error) {
	category := req.ToDomain()
	if req.ImageID != nil {
		// Assuming we handle image association in repository or here
		// For now, let's just pass the ID if the domain supports it, or we need to fetch image
		// The domain model has Image *Image, but ToModel uses ImageID.
		// Let's assume the repository handles the ImageID if it's set in the model.
		// But req.ToDomain() doesn't set ImageID because domain.Category doesn't have it directly, it has Image struct.
		// We might need to adjust this.
		// Let's check domain.Category again.
		// It has Image *Image.
		// So we might need to fetch the image first or just set the ID in the model in the repository.
		// But the repository takes domain.Category.
		// Let's look at CategoryRepository.CreateCategory.
		// It calls category.ToModel().
		// category.ToModel() uses c.Image.ID.
		// So we need to set c.Image with an Image struct that has the ID.
		category.Image = &domain.Image{ID: *req.ImageID}
	}

	if s.categoryRepository.IsExist(ctx, category.Name, category.Slug) {
		return nil, common.ErrConflict(ctx, "Category", "Name or slug already exists")
	}

	if err := s.categoryRepository.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return dto.NewCategoryResponse(category), nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*dto.CategoryResponse, *common.Error) {
	category, err := s.categoryRepository.GetCategoryDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.NewCategoryResponse(category), nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*dto.CategoryResponse, *common.Error) {
	categories, err := s.categoryRepository.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var categoryResponses []*dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, dto.NewCategoryResponse(category))
	}
	return categoryResponses, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, *common.Error) {
	category, err := s.categoryRepository.GetCategoryByID(ctx, req.ID)
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
		category.Image = &domain.Image{ID: *req.ImageID}
	}

	if err := s.categoryRepository.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	// Fetch updated details (e.g. image)
	updatedCategory, err := s.categoryRepository.GetCategoryDetail(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewCategoryResponse(updatedCategory), nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) *common.Error {
	return s.categoryRepository.DeleteCategory(ctx, id)
}

func (s *CategoryService) GetProductsByCategory(ctx context.Context, categoryID uint, page, limit int) (*dto.PaginationResponse[dto.ProductResponse], *common.Error) {
	offset := (page - 1) * limit
	products, total, err := s.productRepository.GetProductsByCategoryID(ctx, categoryID, offset, limit)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, *dto.NewProductResponse(p))
	}

	paginationReq := dto.PaginationRequest{Page: page, Size: limit}
	response := dto.NewPaginationResponse(productResponses, total, paginationReq)
	return &response, nil
}
