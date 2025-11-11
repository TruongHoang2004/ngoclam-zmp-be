package model

import (
	"time"
)

type Image struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"type:varchar(255);not null" json:"url"`
	Hash      string    `gorm:"type:varchar(255);not null" json:"hash"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Image) TableName() string {
	return "images"
}
