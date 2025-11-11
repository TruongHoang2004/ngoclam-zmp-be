package repository

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db          *gorm.DB
	imageKitURL string
	publicKey   string
	privateKey  string
	folder      string
}

func NewImageRepository(db *gorm.DB, config *config.Config) *ImageRepository {
	return &ImageRepository{
		db:          db,
		publicKey:   config.ImageKitPublicKey,
		privateKey:  config.ImageKitPrivateKey,
		folder:      "NgocLamZMP",
		imageKitURL: config.ImageKitEndpoint,
	}
}

// Create (Upload) image
func (r *ImageRepository) UploadImage(ctx context.Context, fileName string, fileData []byte) (*model.Image, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Encode file to base64
	b64 := base64.StdEncoding.EncodeToString(fileData)

	_ = writer.WriteField("file", "data:image/jpeg;base64,"+b64)
	_ = writer.WriteField("fileName", fileName)
	if r.folder != "" {
		_ = writer.WriteField("folder", r.folder)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", r.imageKitURL, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.publicKey, r.privateKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("upload failed: %s", string(data))
	}

	var resp struct {
		URL    string `json:"url"`
		FileID string `json:"fileId"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	img := &model.Image{
		URL:  resp.URL,
		Hash: resp.FileID,
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

// Update (replace image)
func (r *ImageRepository) UpdateImage(ctx context.Context, id uint, fileName string, fileData []byte) (*model.Image, error) {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete old image from ImageKit
	if img.Hash != "" {
		delURL := fmt.Sprintf("https://api.imagekit.io/v1/files/%s", img.Hash)
		req, _ := http.NewRequestWithContext(ctx, "DELETE", delURL, nil)
		req.SetBasicAuth(r.publicKey, r.privateKey)
		client := &http.Client{}
		_, _ = client.Do(req)
	}

	// Upload new image
	return r.UploadImage(ctx, fileName, fileData)
}

// Delete
func (r *ImageRepository) DeleteImage(ctx context.Context, id uint) error {
	img, err := r.GetImageByID(ctx, id)
	if err != nil {
		return err
	}

	if img.Hash != "" {
		delURL := fmt.Sprintf("https://api.imagekit.io/v1/files/%s", img.Hash)
		req, _ := http.NewRequestWithContext(ctx, "DELETE", delURL, nil)
		req.SetBasicAuth(r.publicKey, r.privateKey)
		client := &http.Client{}
		_, _ = client.Do(req)
	}

	return r.db.WithContext(ctx).Delete(&model.Image{}, id).Error
}
