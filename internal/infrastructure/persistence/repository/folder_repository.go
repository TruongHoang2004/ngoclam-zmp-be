package repository

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) CreateFolder(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *FolderRepository) GetFolderByID(ctx context.Context, id uint) (*model.Folder, error) {
	var f model.Folder
	if err := r.db.WithContext(ctx).First(&f, id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FolderRepository) ListFolders(ctx context.Context, offset int, limit int) ([]*model.Folder, int64, error) {
	query := r.db.Model(&model.Folder{})

	var total int64
	query.Count(&total)

	var list []*model.Folder
	if err := query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *FolderRepository) UpdateFolder(ctx context.Context, folder *model.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

func (r *FolderRepository) DeleteFolder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Folder{}, id).Error
}
