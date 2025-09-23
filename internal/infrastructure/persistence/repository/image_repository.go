package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type ImageRepositoryImpl struct {
	db *gorm.DB
	ik *imagekit.ImageKit
}

func NewImageRepository(db *gorm.DB, cfg *config.Config) entity.ImageRepository {
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  cfg.ImageKitPrivateKey,
		PublicKey:   cfg.ImageKitPublicKey,
		UrlEndpoint: cfg.ImageKitEndpoint,
	})

	return &ImageRepositoryImpl{db: db, ik: ik}
}

// SaveFile upload file lên ImageKit và insert DB record
func (r *ImageRepositoryImpl) SaveFile(ctx context.Context, file *multipart.FileHeader) (*entity.Image, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer src.Close()

	// detect content-type
	buf := make([]byte, 512)
	n, _ := src.Read(buf)
	ct := http.DetectContentType(buf[:n])

	allowed := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}
	if !allowed[ct] {
		return nil, fmt.Errorf("unsupported content type: %s", ct)
	}

	// reset về đầu để hash
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, err
	}
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// kiểm tra trùng lặp trong DB
	var existing model.Image
	if err := r.db.WithContext(ctx).Where("hash = ?", fileHash).First(&existing).Error; err == nil {
		return existing.ToDomain(), nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// reset về đầu để upload
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// thêm extension nếu file không có
	exts, _ := mime.ExtensionsByType(ct)
	filename := file.Filename
	if len(exts) > 0 && filepathExt(filename) == "" {
		filename += exts[0]
	}

	// upload lên ImageKit bằng io.Reader
	trueValue := true
	uploadRes, err := r.ik.Uploader.Upload(ctx, src, uploader.UploadParam{
		FileName:          filename,
		UseUniqueFileName: &trueValue,
	})
	if err != nil {
		return nil, fmt.Errorf("upload to imagekit failed: %w", err)
	}

	// lưu metadata vào DB
	img := &model.Image{
		URL:      uploadRes.Data.Url,    // URL public của ảnh
		Hash:     fileHash,              // checksum để chống trùng
		IKFileID: uploadRes.Data.FileId, // để xóa ảnh sau này
	}
	if err := r.db.WithContext(ctx).Create(img).Error; err != nil {
		// rollback: xóa trên imagekit nếu DB fail
		_, _ = r.ik.Media.DeleteFile(context.Background(), uploadRes.Data.FileId)
		return nil, err
	}

	return img.ToDomain(), nil
}

// SetImageLocation chèn hoặc cập nhật ImagePlacement, tự động tăng DisplayOrder
func (r *ImageRepositoryImpl) SetImageLocation(ctx context.Context, imageID uint, location string) error {
	var maxOrder int

	// Lấy max(display_order) hiện tại
	err := r.db.WithContext(ctx).
		Model(&model.ImagePlacement{}).
		Where("location = ?", location).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	if err != nil {
		return err
	}

	imgPlacement := &model.ImagePlacement{
		ImageID:      imageID,
		Location:     location,
		DisplayOrder: maxOrder + 1,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	// ON CONFLICT: cần unique constraint trên (image_id, location)
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "image_id"}, {Name: "location"}},
		UpdateAll: true,
	}).Create(imgPlacement).Error; err != nil {
		return err
	}

	return nil
}

func (r *ImageRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Image, error) {
	var img model.Image
	if err := r.db.WithContext(ctx).First(&img, id).Error; err != nil {
		return nil, err
	}
	return img.ToDomain(), nil
}

// FindByPlacement implements entity.ImageRepository.
func (r *ImageRepositoryImpl) FindByPlacement(ctx context.Context, location string) ([]*entity.Image, error) {
	var imgs []model.Image
	if err := r.db.WithContext(ctx).
		Joins("JOIN image_placements ON images.id = image_placements.image_id").
		Where("image_placements.location = ?", location).
		Find(&imgs).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Image, len(imgs))
	for i, img := range imgs {
		result[i] = img.ToDomain()
	}
	return result, nil
}

// FindAll implements entity.ImageRepository.
func (r *ImageRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Image, error) {
	var imgs []model.Image
	if err := r.db.WithContext(ctx).Find(&imgs).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Image, len(imgs))
	for i, img := range imgs {
		result[i] = img.ToDomain()
	}
	return result, nil
}

func (r *ImageRepositoryImpl) Delete(ctx context.Context, id uint) error {
	var img model.Image
	if err := r.db.WithContext(ctx).First(&img, id).Error; err != nil {
		return err
	}

	// xóa trên ImageKit
	if img.IKFileID != "" {
		_, err := r.ik.Media.DeleteFile(ctx, img.IKFileID)
		if err != nil {
			log.Printf("warn: failed to delete imagekit file %s: %v", img.IKFileID, err)
		}
	}

	return r.db.Delete(&img).Error
}

// helper để lấy extension từ filename
func filepathExt(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i:]
		}
	}
	return ""
}
