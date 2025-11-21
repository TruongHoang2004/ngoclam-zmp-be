package domain

import (
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type Folder struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	ParentID *uint     `json:"parent_id,omitempty"`
	Parent   *Folder   `json:"parent,omitempty"`
	Children *[]Folder `json:"children,omitempty"`

	Images *[]Image `json:"images,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFolderDomain(folder *model.Folder) *Folder {
	return &Folder{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		ParentID:    folder.ParentID,
		Parent:      nil,
		Children:    nil,
		Images:      nil,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}
}

func (f *Folder) AddParent(parent *Folder) {
	f.Parent = parent
}

func (f *Folder) AddChildren(children []Folder) {
	f.Children = &children
}

func (f *Folder) AddImages(images []Image) {
	f.Images = &images
}
