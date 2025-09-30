package application

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
)

type ImageService struct {
	imageRepo entity.ImageRepository
}

func NewImageService(imageRepo entity.ImageRepository) *ImageService {
	return &ImageService{
		imageRepo: imageRepo,
	}
}

func (s *ImageService) SaveImageByURL(ctx context.Context, url string) (*entity.Image, error) {
	image, err := s.imageRepo.SaveByURL(ctx, url)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_images_hash\" (SQLSTATE 23505)" {
			return nil, NewConflictError("image already exists")
		}
		return nil, NewInternalServerError(fmt.Sprintf("cannot save image by URL: %v", err))
	}
	return image, nil
}

func (s *ImageService) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (*entity.Image, error) {
	// Validate file size (e.g., max 5MB)
	const maxFileSize = 20 << 20 // 20MB
	if fileHeader.Size > maxFileSize {
		return nil, NewUnsupportedMediaTypeError(fmt.Sprintf("file size exceeds the limit of %d bytes", maxFileSize))
	}
	image, err := s.imageRepo.SaveFile(ctx, fileHeader)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"idx_images_hash\" (SQLSTATE 23505)" {
			return nil, NewConflictError("image already exists")
		}
		return nil, NewInternalServerError(fmt.Sprintf("cannot save image: %v", err))
	}
	return image, nil
}

func (s *ImageService) GetImageByID(ctx context.Context, id uint) (*entity.Image, error) {
	image, err := s.imageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot find image by ID: %w", err)
	}
	return image, nil
}

func (s *ImageService) GetImageByLocation(ctx context.Context, location string) ([]*entity.Image, error) {
	images, err := s.imageRepo.FindByPlacement(ctx, location)
	if err != nil {
		return nil, NewNotFoundError(fmt.Sprintf("cannot find images by location: %v", err))
	}
	return images, nil
}

func (s *ImageService) SetImageLocation(ctx context.Context, imageID uint, location string) error {
	if err := s.imageRepo.SetImageLocation(ctx, imageID, location); err != nil {
		return fmt.Errorf("cannot set image location: %w", err)
	}
	return nil
}

func (s *ImageService) ListImages(ctx context.Context) ([]*entity.Image, error) {
	images, err := s.imageRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot list images: %w", err)
	}
	return images, nil
}

func (s *ImageService) DeleteImage(ctx context.Context, id uint) error {
	if err := s.imageRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("cannot delete image: %w", err)
	}
	return nil
}
