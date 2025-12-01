package model

import "time"

type Folder struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_folder_parent_name" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	ParentID *uint   `gorm:"uniqueIndex:idx_folder_parent_name;index" json:"parent_id,omitempty"`
	Parent   *Folder `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	Children []Folder `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Images   []Image  `gorm:"foreignKey:FolderID" json:"images,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Folder) TableName() string {
	return "folders"
}
