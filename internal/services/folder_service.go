package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type FolderService struct {
	repo *repositories.FolderRepository
}

func NewFolderService(repo *repositories.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(ctx context.Context, name string, description string) (*domain.Folder, *common.Error) {
	f := &domain.Folder{
		Name:        name,
		Description: description,
	}
	if err := s.repo.CreateFolder(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FolderService) GetFolderByID(ctx context.Context, id uint) (*domain.Folder, *common.Error) {
	return s.repo.GetFolderByID(ctx, id)
}

func (s *FolderService) ListFolders(ctx context.Context, page int, size int) ([]*domain.Folder, int64, *common.Error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	return s.repo.ListFolders(ctx, offset, size)
}

func (s *FolderService) UpdateFolder(ctx context.Context, folderId uint, req *dto.UpdateFolderRequest) (*domain.Folder, *common.Error) {
	folder, err := s.repo.GetFolderByID(ctx, folderId)
	if err != nil {
		return nil, err
	}
	folder.Name = req.Name
	folder.Description = req.Description
	if updated, err := s.repo.UpdateFolder(ctx, folder); err != nil {
		return nil, err
	} else {
		folder = updated
	}
	return folder, nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, id uint) error {
	return s.repo.DeleteFolder(ctx, id)
}
