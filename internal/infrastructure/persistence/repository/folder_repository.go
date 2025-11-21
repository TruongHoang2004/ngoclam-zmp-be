package repository

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) CreateFolder(ctx context.Context, folder *domain.Folder) error {
	f := &model.Folder{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
	}

	return r.db.WithContext(ctx).Create(f).Error
}

func (r *FolderRepository) GetFolderByID(ctx context.Context, id uint) (*domain.Folder, error) {
	var folder *model.Folder
	if err := r.db.WithContext(ctx).First(&folder, id).Error; err != nil {
		return nil, err
	}
	return domain.NewFolderDomain(folder), nil
}

func (r *FolderRepository) ListFolders(ctx context.Context, offset int, limit int) ([]*domain.Folder, int64, error) {
	query := r.db.Model(&model.Folder{})

	var total int64
	query.Count(&total)

	var list []*model.Folder
	if err := query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	domainList := make([]*domain.Folder, 0, len(list))
	for _, f := range list {
		domainList = append(domainList, domain.NewFolderDomain(f))
	}
	return domainList, total, nil
}

func (r *FolderRepository) UpdateFolder(ctx context.Context, folder *domain.Folder) error {
	return r.db.WithContext(ctx).Save(&model.Folder{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		ParentID:    folder.ParentID,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}).Error
}

func (r *FolderRepository) DeleteFolder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Folder{}, id).Error
}
