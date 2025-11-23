package repositories

import (
	"context"
	"errors"
	"io"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/sdk/imagekit"
	"gorm.io/gorm"
)

type ImageRepository struct {
	*baseRepository
	imageKitClient *imagekit.ImageKitClient
	folder         string
}

func NewImageRepository(base *baseRepository, imageKitClient *imagekit.ImageKitClient) *ImageRepository {
	return &ImageRepository{
		baseRepository: base,
		imageKitClient: imageKitClient,
		folder:         "NgocLamZMP",
	}
}

// UploadImage uploads an image from byte data
func (r *ImageRepository) UploadImage(ctx context.Context, fileName string, fileData []byte) (*domain.Image, *common.Error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadImageFromBytes(ctx, fileData, fileName, opts)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}

	return domain.NewImageDomain(img), nil
}

// UploadImageFromReader uploads an image from io.Reader
func (r *ImageRepository) UploadImageFromReader(ctx context.Context, file io.Reader, fileName string) (*domain.Image, *common.Error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadImage(ctx, file, fileName, opts)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}

	return domain.NewImageDomain(img), nil
}

// UploadImageFromURL uploads an image from a URL
func (r *ImageRepository) UploadImageFromURL(ctx context.Context, url string, fileName string) (*domain.Image, *common.Error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadFromURL(ctx, url, fileName, opts)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}

	return domain.NewImageDomain(img), nil
}

// Read
func (r *ImageRepository) GetImageByID(ctx context.Context, id uint) (*model.Image, *common.Error) {
	var img model.Image
	if err := r.db.WithContext(ctx).First(&img, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "Image", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}
	return &img, nil
}

func (r *ImageRepository) GetAllImages(ctx context.Context, page int, limit int) ([]*domain.Image, int64, *common.Error) {
	var images []*domain.Image
	offset := (page - 1) * limit

	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Image{}).Count(&total).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&images).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return images, total, nil
}

// Update (replace image) updates an image from byte data
func (r *ImageRepository) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*domain.Image, *common.Error) {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete old image from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Upload new image
	return r.UploadImage(ctx, fileName, fileData)
}

// UpdateImageFromReader updates an image from io.Reader
func (r *ImageRepository) UpdateImageFromReader(ctx context.Context, id uint, file io.Reader, fileName string) (*domain.Image, *common.Error) {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete old image from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Upload new image
	return r.UploadImageFromReader(ctx, file, fileName)
}

// UpdateImageFromURL updates an image from a URL
func (r *ImageRepository) UpdateImageFromURL(ctx context.Context, id uint, url string, fileName string) (*domain.Image, *common.Error) {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete old image from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Upload new image
	image, err := r.UploadImageFromURL(ctx, url, fileName)
	return image, err
}

// Delete deletes an image from database and ImageKit
func (r *ImageRepository) DeleteImage(ctx context.Context, id uint) *common.Error {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Delete from database
	return r.returnError(ctx, r.db.WithContext(ctx).Delete(&model.Image{}, id).Error)
}
