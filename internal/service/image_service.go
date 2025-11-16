package service

import (
	"context"
	"io"

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

// UploadImage uploads an image from byte data
func (s *ImageService) UploadImage(ctx context.Context, fileName string, fileData []byte) (*model.Image, error) {
	return s.imageRepository.UploadImage(ctx, fileName, fileData)
}

// UploadImageFromReader uploads an image from io.Reader
func (s *ImageService) UploadImageFromReader(ctx context.Context, file io.Reader, fileName string) (*model.Image, error) {
	return s.imageRepository.UploadImageFromReader(ctx, file, fileName)
}

// UploadImageFromURL uploads an image from a URL
func (s *ImageService) UploadImageFromURL(ctx context.Context, url string, fileName string) (*model.Image, error) {
	return s.imageRepository.UploadImageFromURL(ctx, url, fileName)
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

// UpdateImage updates an image from byte data
func (s *ImageService) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*model.Image, error) {
	image, err := s.imageRepository.UpdateImage(ctx, id, fileName, fileData)
	return image, err
}

// UpdateImageFromReader updates an image from io.Reader
func (s *ImageService) UpdateImageFromReader(ctx context.Context, id uint, file io.Reader, fileName string) (*model.Image, error) {
	return s.imageRepository.UpdateImageFromReader(ctx, id, file, fileName)
}

// UpdateImageFromURL updates an image from a URL
func (s *ImageService) UpdateImageFromURL(ctx context.Context, id uint, url string, fileName string) (*model.Image, error) {
	return s.imageRepository.UpdateImageFromURL(ctx, id, url, fileName)
}

func (s *ImageService) DeleteImage(ctx context.Context, id uint) error {
	return s.imageRepository.DeleteImage(ctx, id)
}
