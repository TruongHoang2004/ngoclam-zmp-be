package repositories

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	*baseRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		baseRepository: NewBaseRepository(db),
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) *common.Error {
	if err := r.db.Create(order).Error; err != nil {
		return common.ErrSystemError(ctx, err.Error())
	}
	return nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id uint) (*model.Order, *common.Error) {
	var order model.Order
	if err := r.db.Preload("OrderItems").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrNotFound(ctx, "Order", "not found")
		}
		return nil, common.ErrSystemError(ctx, err.Error())
	}
	return &order, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context, offset int, limit int) ([]*model.Order, int64, *common.Error) {
	var orders []*model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, common.ErrSystemError(ctx, err.Error())
	}

	if err := r.db.Preload("OrderItems").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		return nil, 0, common.ErrSystemError(ctx, err.Error())
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *model.Order) *common.Error {
	if err := r.db.Save(order).Error; err != nil {
		return common.ErrSystemError(ctx, err.Error())
	}
	return nil
}
