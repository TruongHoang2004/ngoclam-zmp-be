package model

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain/entity"

type ImageModel struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	URL      string `gorm:"not null;size:512"`             // link public của ảnh trên ImageKit
	IKFileID string `gorm:"not null;size:128;uniqueIndex"` // fileId trong ImageKit
	Hash     string `gorm:"size:64;uniqueIndex;not null"`  // sha256 checksum
}

func (ImageModel) TableName() string {
	return "images"
}

// Map sang domain entity
func (m *ImageModel) ToDomain() *entity.Image {
	return &entity.Image{
		ID:       m.ID,
		URL:      m.URL,
		IKFileID: m.IKFileID,
		Hash:     m.Hash,
	}
}

func MapImageToModel(eimg *entity.Image) *ImageModel {
	if eimg == nil {
		return nil
	}
	return &ImageModel{
		ID:       eimg.ID,
		URL:      eimg.URL,
		IKFileID: eimg.IKFileID,
		Hash:     eimg.Hash,
	}
}
