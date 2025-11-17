package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"

type FolderResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateFolderRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

func NewFolderResponse(f *model.Folder) *FolderResponse {
	return &FolderResponse{
		ID:          f.ID,
		Name:        f.Name,
		Description: f.Description,
	}
}
