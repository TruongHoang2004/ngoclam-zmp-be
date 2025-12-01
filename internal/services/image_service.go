package services

import (
	"context"
	"io"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
)

type ImageService struct {
	*baseService
	imageRepository *repositories.ImageRepository
}

func NewImageService(base *baseService, imageRepo *repositories.ImageRepository) *ImageService {
	return &ImageService{
		baseService:     base,
		imageRepository: imageRepo,
	}
}

// UploadImage uploads an image from byte data
func (s *ImageService) UploadImage(ctx context.Context, fileName string, fileData []byte) (*model.Image, *common.Error) {
	return s.imageRepository.UploadImage(ctx, fileName, fileData)
}

// UploadImageFromReader uploads an image from io.Reader
func (s *ImageService) UploadImageFromReader(ctx context.Context, file io.Reader, fileName string) (*model.Image, *common.Error) {
	return s.imageRepository.UploadImageFromReader(ctx, file, fileName)
}

// UploadImageFromURL uploads an image from a URL
func (s *ImageService) UploadImageFromURL(ctx context.Context, url string, fileName string) (*model.Image, *common.Error) {
	return s.imageRepository.UploadImageFromURL(ctx, url, fileName)
}

func (s *ImageService) GetImageByID(ctx context.Context, id uint) (*model.Image, *common.Error) {
	img, err := s.imageRepository.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, common.ErrNotFound(ctx, "Image", "not found").SetSource(common.CurrentService)
	}
	return img, nil
}

func (s *ImageService) GetAllImages(ctx context.Context, page int, limit int) ([]*model.Image, int64, *common.Error) {
	list, total, err := s.imageRepository.GetAllImages(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateImage updates an image from byte data
func (s *ImageService) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*model.Image, *common.Error) {
	image, err := s.imageRepository.UpdateImage(ctx, id, fileName, fileData)
	if err != nil {
		return nil, err
	}
	return image, nil
}

// UpdateImageFromReader updates an image from io.Reader
func (s *ImageService) UpdateImageFromReader(ctx context.Context, id uint, file io.Reader, fileName string) (*model.Image, *common.Error) {
	return s.imageRepository.UpdateImageFromReader(ctx, id, file, fileName)
}

// UpdateImageFromURL updates an image from a URL
func (s *ImageService) UpdateImageFromURL(ctx context.Context, id uint, url string, fileName string) (*model.Image, *common.Error) {
	return s.imageRepository.UpdateImageFromURL(ctx, id, url, fileName)
}

func (s *ImageService) DeleteImage(ctx context.Context, id uint) *common.Error {
	return s.imageRepository.DeleteImage(ctx, id)
}
