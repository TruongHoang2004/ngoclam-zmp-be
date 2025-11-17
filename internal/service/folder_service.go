package service

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
)

type FolderService struct {
	repo *repository.FolderRepository
}

func NewFolderService(repo *repository.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(ctx context.Context, name string, description string) (*model.Folder, error) {
	f := &model.Folder{
		Name:        name,
		Description: description,
	}
	if err := s.repo.CreateFolder(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FolderService) GetFolderByID(ctx context.Context, id uint) (*model.Folder, error) {
	return s.repo.GetFolderByID(ctx, id)
}

func (s *FolderService) ListFolders(ctx context.Context, page int, size int) ([]*model.Folder, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size
	return s.repo.ListFolders(ctx, offset, size)
}

func (s *FolderService) UpdateFolder(ctx context.Context, folder *model.Folder) (*model.Folder, error) {
	if err := s.repo.UpdateFolder(ctx, folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (s *FolderService) DeleteFolder(ctx context.Context, id uint) error {
	return s.repo.DeleteFolder(ctx, id)
}
