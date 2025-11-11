package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"

type ImageResponse struct {
	ID   uint   `json:"id"`
	URL  string `json:"url"`
	Hash string `json:"hash"`
}

type UploadImageRequest struct {
	FileName string `json:"fileName"`
	FileData []byte `json:"fileData"`
}

func NewImageResponse(image *model.Image) *ImageResponse {
	return &ImageResponse{
		ID:   image.ID,
		URL:  image.URL,
		Hash: image.Hash,
	}
}
