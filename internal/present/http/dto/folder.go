package dto

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
)

type FolderResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewFolderResponse(f *model.Folder) *FolderResponse {
	return &FolderResponse{
		ID:          f.ID,
		Name:        f.Name,
		Description: f.Description,
	}
}

type CreateFolderRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

type UpdateFolderRequest struct {
	Name        *string `json:"name" `
	Description *string `json:"description,omitempty"`
	ParentID    *uint   `json:"parent_id,omitempty"`
}

func (r *CreateFolderRequest) ToModel() *model.Folder {
	return &model.Folder{
		Name:        r.Name,
		Description: r.Description,
		ParentID:    r.ParentID,
	}
}
