package repository

import (
	"context"
	"io"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/sdk/imagekit"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db             *gorm.DB
	imageKitClient *imagekit.ImageKitClient
	folder         string
}

func NewImageRepository(db *gorm.DB, imageKitClient *imagekit.ImageKitClient) *ImageRepository {
	return &ImageRepository{
		db:             db,
		imageKitClient: imageKitClient,
		folder:         "NgocLamZMP",
	}
}

// UploadImage uploads an image from byte data
func (r *ImageRepository) UploadImage(ctx context.Context, fileName string, fileData []byte) (*model.Image, error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadImageFromBytes(ctx, fileData, fileName, opts)
	if err != nil {
		return nil, err
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, err
	}

	return img, nil
}

// UploadImageFromReader uploads an image from io.Reader
func (r *ImageRepository) UploadImageFromReader(ctx context.Context, file io.Reader, fileName string) (*model.Image, error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadImage(ctx, file, fileName, opts)
	if err != nil {
		return nil, err
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, err
	}

	return img, nil
}

// UploadImageFromURL uploads an image from a URL
func (r *ImageRepository) UploadImageFromURL(ctx context.Context, url string, fileName string) (*model.Image, error) {
	opts := &imagekit.UploadOptions{
		Folder:            r.folder,
		UseUniqueFileName: true,
	}

	result, err := r.imageKitClient.UploadFromURL(ctx, url, fileName, opts)
	if err != nil {
		return nil, err
	}

	img := &model.Image{
		URL:  result.Url,
		Hash: result.FileId,
	}

	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		return nil, err
	}

	return img, nil
}

// Read
func (r *ImageRepository) GetImageByID(ctx context.Context, id uint) (*model.Image, error) {
	var img model.Image
	if err := r.db.WithContext(ctx).First(&img, id).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepository) GetAllImages(ctx context.Context, page int, limit int) ([]*model.Image, error) {
	var images []*model.Image
	offset := page * limit
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

// Update (replace image) updates an image from byte data
func (r *ImageRepository) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*model.Image, error) {
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
func (r *ImageRepository) UpdateImageFromReader(ctx context.Context, id uint, file io.Reader, fileName string) (*model.Image, error) {
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
func (r *ImageRepository) UpdateImageFromURL(ctx context.Context, id uint, url string, fileName string) (*model.Image, error) {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete old image from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Upload new image
	return r.UploadImageFromURL(ctx, url, fileName)
}

// Delete deletes an image from database and ImageKit
func (r *ImageRepository) DeleteImage(ctx context.Context, id uint) error {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from ImageKit
	if img.Hash != "" {
		_ = r.imageKitClient.DeleteFile(ctx, img.Hash)
	}

	// Delete from database
	return r.db.WithContext(ctx).Delete(&model.Image{}, id).Error
}
