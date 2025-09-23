package model

type ImagePlacement struct {
	ID           uint   `gorm:"column:id;primaryKey;autoIncrement"`
	ImageID      uint   `gorm:"uniqueIndex:idx_image_location,priority:1"`
	Location     string `gorm:"uniqueIndex:idx_image_location,priority:2"`
	DisplayOrder int    `gorm:"column:display_order;default:0"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoUpdateTime"`
}

func (ImagePlacement) TableName() string {
	return "image_placements"
}

func NewImagePlacement(imageID uint, location string, width int, height int) *ImagePlacement {
	return &ImagePlacement{
		ImageID:      imageID,
		Location:     location,
		DisplayOrder: 0,
	}
}
