package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type FolderService struct {
	repo *repositories.FolderRepository
}

func NewFolderService(repo *repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(ctx context.Context, req *dto.CreateFolderRequest) *common.Error {
	f := req.ToModel()
	if req.ParentID != nil {
		f.ParentID = req.ParentID
	}
	return s.repo.CreateFolder(ctx, f)
}

func (s *FolderService) GetFolderByID(ctx context.Context, id uint) (*model.Folder, *common.Error) {
	return s.repo.GetFolderByID(ctx, id)
}

func (s *FolderService) ListFolders(ctx context.Context, page int, size int) ([]*model.Folder, int64, *common.Error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	return s.repo.ListFolders(ctx, offset, size)
}

func (s *FolderService) UpdateFolder(ctx context.Context, id uint, req *dto.UpdateFolderRequest) (*model.Folder, *common.Error) {
	folder, err := s.repo.GetFolderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		folder.Name = *req.Name
	}
	if req.Description != nil {
		folder.Description = *req.Description
	}
	if req.ParentID != nil {
		folder.ParentID = req.ParentID
	}

	return s.repo.UpdateFolder(ctx, folder)
}

func (s *FolderService) DeleteFolder(ctx context.Context, id uint) error {
	return s.repo.DeleteFolder(ctx, id)
}
