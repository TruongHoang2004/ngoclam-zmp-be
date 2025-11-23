package repositories

import (
	"context"
	"errors"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/domain"
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

func (r *FolderRepository) CreateFolder(ctx context.Context, folder *domain.Folder) *common.Error {
	f := &model.Folder{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
	}

	return r.returnError(ctx, r.db.WithContext(ctx).Create(f).Error)
}

func (r *FolderRepository) GetFolderByID(ctx context.Context, id uint) (*domain.Folder, *common.Error) {
	folder := &model.Folder{}
	cond := clause.Eq{Column: "id", Value: id}

	if err := r.db.Clauses(cond).Take(folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "folder", "not found").SetSource(common.CurrentService)
		}
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
