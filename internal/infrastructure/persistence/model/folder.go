package model

import "time"

// Folder represents a nested folder that can contain images and other folders.
type Folder struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Name is unique among siblings (same parent). We use a composite unique index
	// idx_folder_parent_name to enforce uniqueness of (parent_id, name).
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_folder_parent_name" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Parent relationship for nested folders. ParentID is nullable for root folders.
	ParentID *uint     `gorm:"uniqueIndex:idx_folder_parent_name;index" json:"parent_id,omitempty"`
	Parent   *Folder   `json:"parent,omitempty"`
	Children []*Folder `json:"children,omitempty"`

	// Images contained in this folder
	Images []*Image `gorm:"foreignKey:FolderID" json:"images,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Folder) TableName() string {
	return "folders"
}
