package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	orderRepository   *repositories.OrderRepository
	productRepository *repositories.ProductRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository) *OrderService {
	return &OrderService{
		orderRepository:   orderRepo,
		productRepository: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*model.Order, *common.Error) {
	// 1. Validate items and Calculate total
	var orderItems []model.OrderItem
	totalAmount := decimal.Zero

	for _, itemReq := range req.Items {
		// Get Product details for snapshot
		product, err := s.productRepository.GetProductByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, common.ErrNotFound(ctx, "Product", "not found")
		}

		// Create Snapshot
		snapshot := &model.ProductSnapshot{
			ProductID:   product.ID,
			Name:        product.Name,
			Price:       decimal.NewFromInt(product.Price),
			Description: "", // Optional
		}
		if product.Description != nil {
			snapshot.Description = *product.Description
		}
		// Assuming first image as main image for snapshot, or existing logic
		if len(product.ProductImages) > 0 {
			// This logic depends on how images are loaded.
			// ProductRepository.GetProductByID might need to preload images if we want them in snapshot.
			// Let's assume for now we might need to fetch images if not preloaded, or just skip if complex.
			// Checking ProductRepository.GetProductByID implementation:
			// It calls GetProductDetailByID which likely preloads.
			// Let's check the product model again. ProductImages []ProductImage.
			for _, img := range product.ProductImages {
				if img.IsMain && img.Image != nil {
					snapshot.ImageURL = img.Image.URL // Assuming Image model has Url
					break
				}
			}
			// Fallback if no main image found but images exist
			if snapshot.ImageURL == "" && len(product.ProductImages) > 0 && product.ProductImages[0].Image != nil {
				snapshot.ImageURL = product.ProductImages[0].Image.URL
			}
		}

		price := decimal.NewFromInt(product.Price)
		quantity := decimal.NewFromInt(int64(itemReq.Quantity))
		lineTotal := price.Mul(quantity)
		totalAmount = totalAmount.Add(lineTotal)

		orderItems = append(orderItems, model.OrderItem{
			ProductSnapshot: snapshot,
			Quantity:        itemReq.Quantity,
			Price:           price,
		})
	}

	// 2. Create Order Model
	custInfo := &model.CustomerInfo{
		Name:    req.CustomerInfo.Name,
		Phone:   req.CustomerInfo.Phone,
		Address: req.CustomerInfo.Address,
	}

	order := &model.Order{
		CustomerInfo: custInfo,
		TotalAmount:  totalAmount,
		Status:       "pending",
		OrderItems:   orderItems,
	}

	// 3. Save to DB
	if err := s.orderRepository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context, page int, size int) ([]*model.Order, int64, *common.Error) {
	offset := (page - 1) * size
	return s.orderRepository.ListOrders(ctx, offset, size)
}

func (s *OrderService) GetOrder(ctx context.Context, id uint) (*model.Order, *common.Error) {
	return s.orderRepository.GetOrder(ctx, id)
}
