package model

type ImagePlacement struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement"`
	ImageID   uint   `gorm:"column:image_id;not null"`
	Location  string `gorm:"column:location;not null"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime"`
}

func (ImagePlacement) TableName() string {
	return "image_placements"
}

func NewImagePlacement(imageID uint, location string, width int, height int) *ImagePlacement {
	return &ImagePlacement{
		ImageID:  imageID,
		Location: location,
	}
}
