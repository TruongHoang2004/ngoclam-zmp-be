package repositories

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *baseRepository {
	return &baseRepository{
		db: db,
	}
}

func (b *baseRepository) returnError(ctx context.Context, err error) *common.Error {
	return common.ErrSystemError(ctx, err.Error()).SetSource(common.CurrentService)
}

func (b *baseRepository) ApplyFilter(db *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		db = db.Where(key, value)
	}
	return db
}
