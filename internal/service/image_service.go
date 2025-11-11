package service

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
)

type ImageService struct {
	imageRepository *repository.ImageRepository
}

func NewImageService(imageRepo *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepo,
	}
}

func (s *ImageService) UploadImage(ctx context.Context, fileName string, fileData []byte) (*model.Image, error) {
	return s.imageRepository.UploadImage(ctx, fileName, fileData)
}

func (s *ImageService) GetImageByID(ctx context.Context, id uint) (*model.Image, error) {

	image, err := s.imageRepository.GetImageByID(ctx, id)

	if image == nil {
		return nil, common.NotFound("Image id not found")
	}

	return image, err
}

func (s *ImageService) GetAllImages(ctx context.Context, page int, limit int) ([]*model.Image, error) {
	list, err := s.imageRepository.GetAllImages(ctx, page, limit)

	return list, err
}

func (s *ImageService) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*model.Image, error) {
	image, err := s.imageRepository.UpdateImage(ctx, id, fileName, fileData)
	return image, err
}

func (s *ImageService) DeleteImage(ctx context.Context, id uint) error {
	return s.imageRepository.DeleteImage(ctx, id)
}
