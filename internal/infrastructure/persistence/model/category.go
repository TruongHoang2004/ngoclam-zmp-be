package model

import (
	"time"
)

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
	Slug string `gorm:"type:varchar(255);not null;uniqueIndex" json:"slug"`

	ImageID   *uint     `gorm:"index" json:"image_id,omitempty"`
	Image     *Image    `gorm:"foreignKey:ImageID" json:"image,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
