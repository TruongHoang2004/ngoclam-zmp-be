package model

type EntityType string

const (
	EntityTypeProduct  EntityType = "product"
	EntityTypeCategory EntityType = "category"
)

type ImageRelated struct {
	ID         uint       `gorm:"column:id;primaryKey;autoIncrement"`
	ImageID    uint       `gorm:"column:image_id;not null"`
	EntityID   uint       `gorm:"column:entity_id;not null"`
	EntityType EntityType `gorm:"column:entity_type;not null"`
	Order      int        `gorm:"column:order;default:0"`

	Image Image `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ImageRelated) TableName() string {
	return "image_related"
}

func CreateImageRelated(imageID uint, entityID uint, entityType EntityType, order int) *ImageRelated {
	return &ImageRelated{
		ImageID:    imageID,
		EntityID:   entityID,
		EntityType: entityType,
		Order:      order,
	}
}
