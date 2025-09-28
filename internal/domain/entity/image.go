package entity

import (
	"context"
	"mime/multipart"
)

// Image represents an image associated with a product
type Image struct {
	ID       uint   `json:"id"`       // Unique identifier for the image
	Path     string `json:"path"`     // URL path to the image
	URL      string `json:"url"`      // URL to access the image
	IKFileID string `json:"ikFileId"` // ImageKit file ID
	Hash     string `json:"hash"`     // Hash of the image
}

// NewImage creates a new Image entity
func NewImage(imagePath string, isPrimary bool) *Image {
	return &Image{
		ID:   0,
		Path: imagePath,
	}
}

type ImageRepository interface {
	SaveByURL(ctx context.Context, url string) (*Image, error)
	SaveFile(ctx context.Context, file *multipart.FileHeader) (*Image, error)
	FindByID(ctx context.Context, id uint) (*Image, error)
	SetImageLocation(ctx context.Context, imageID uint, location string) error
	FindByPlacement(ctx context.Context, location string) ([]*Image, error)
	FindAll(ctx context.Context) ([]*Image, error)
	Delete(ctx context.Context, id uint) error
}
