package domain

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type Image struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Hash      string    `json:"hash"`
	FolderID  *uint     `json:"folder_id,omitempty"`
	Folder    *Folder   `json:"folder,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewImageDomain(image *model.Image) *Image {
	return &Image{
		ID:        image.ID,
		Name:      image.Name,
		URL:       image.URL,
		Hash:      image.Hash,
		FolderID:  image.FolderID,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}

func (i *Image) SetFolder(folder *Folder) {
	i.Folder = folder
}
