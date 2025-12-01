package repositories

import (
	"context"
	"errors"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FolderRepository struct {
	*baseRepository
}

func NewFolderRepository(base *baseRepository) *FolderRepository {
	return &FolderRepository{baseRepository: base}
}

func (r *FolderRepository) CreateFolder(ctx context.Context, folder *model.Folder) *common.Error {
	return r.returnError(ctx, r.db.WithContext(ctx).Create(folder).Error)
}

func (r *FolderRepository) GetFolderByID(ctx context.Context, id uint) (*model.Folder, *common.Error) {
	folder := &model.Folder{}
	cond := clause.Eq{Column: "id", Value: id}

	if err := r.db.Clauses(cond).
		Preload("Parent").
		Preload("Children").
		Preload("Images").
		Take(folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "folder", "not found").SetSource(common.CurrentService)
		}
		return nil, r.returnError(ctx, err)
	}
	return folder, nil
}

func (r *FolderRepository) ListFolders(ctx context.Context, offset int, limit int) ([]*model.Folder, int64, *common.Error) {
	query := r.db.Model(&model.Folder{})

	var total int64
	query.Count(&total)

	var list []*model.Folder
	if err := query.Offset(offset).Limit(limit).
		Preload("Parent").
		Preload("Children").
		Preload("Images").
		Find(&list).Error; err != nil {
		return nil, 0, common.ErrSystemError(ctx, err.Error())
	}
	return list, total, nil
}

func (r *FolderRepository) UpdateFolder(ctx context.Context, folder *model.Folder) (*model.Folder, *common.Error) {
	if err := r.db.WithContext(ctx).Save(folder).Error; err != nil {
		return nil, common.ErrSystemError(ctx, err.Error())
	}
	return folder, nil
}

func (r *FolderRepository) DeleteFolder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Folder{}, id).Error
}
