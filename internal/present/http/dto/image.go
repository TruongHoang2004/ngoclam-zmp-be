package dto

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type ImageResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Hash      string `json:"hash"`
	FolderID  uint   `json:"folder_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UploadImageRequest struct {
	FileName string `json:"fileName"`
	FileData []byte `json:"fileData"`
}

func NewImageResponse(image *model.Image) *ImageResponse {
	var folderID uint
	if image.FolderID != nil {
		folderID = *image.FolderID
	}
	return &ImageResponse{
		ID:        image.ID,
		Name:      image.Name,
		URL:       image.URL,
		Hash:      image.Hash,
		FolderID:  folderID,
		CreatedAt: image.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: image.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
