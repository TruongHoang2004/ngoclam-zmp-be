package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils/casting"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	orderRepository   *repositories.OrderRepository
	productRepository *repositories.ProductRepository
	cfg               *config.Config
}

func NewOrderService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository, cfg *config.Config) *OrderService {
	return &OrderService{
		orderRepository:   orderRepo,
		productRepository: productRepo,
		cfg:               cfg,
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

func (s *OrderService) ProcessZaloCallback(ctx context.Context, req *dto.ZaloCallbackRequest) (*dto.ZaloCallbackResponse, *common.Error) {
	// 1. Verify Message Authentication Code (HMAC-SHA256)
	// data = req.Data
	// mac = hmac256(data, key)
	// assuming key is ZaloAppSecret
	mac := utils.ComputeHmac256(req.Data, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		// return invalid mac
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Data
	// req.Data is a JSON string. We need to parse it to find reference to our order.
	// Common ZaloPay fields: "app_trans_id" or "embed_data".
	// Let's assume we put order ID in "embed_data" or it is derived from "app_trans_id"
	// Example Data: {"app_id": 2553, "app_trans_id": "210608_12345", ...}
	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(req.Data), &dataMap); err != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid data format",
		}, nil
	}

	// Try to get Order ID.
	// Strategy: app_trans_id format YYMMDD_OrderID.
	// Or maybe embed_data contains {"order_id": 123}
	// Let's try to find "embed_data" and parse it as JSON if it's a string, or check map.
	// Simplified assumption for MVP: app_trans_id contains the ID after an underscore.
	// "250101_1" -> Order ID 1
	appTransID, ok := dataMap["app_trans_id"].(string)
	if !ok {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "app_trans_id not found",
		}, nil
	}

	parts := strings.Split(appTransID, "_")
	if len(parts) < 2 {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid app_trans_id format",
		}, nil
	}

	orderIDStr := parts[1]
	orderID, err := casting.StringToUint(orderIDStr)
	if err != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid order id in app_trans_id",
		}, nil
	}

	// 3. Update Order Satus
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "order not found",
		}, nil
	}

	// Check if already paid
	if order.Status == "success" {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    1,
			ReturnMessage: "success",
		}, nil
	}

	order.Status = "success"
	// Optional: store transaction ID from Zalo
	if zpTransID, ok := dataMap["zp_trans_id"].(string); ok {
		order.TransactionID = &zpTransID
	} else if zpTransIDFloat, ok := dataMap["zp_trans_id"].(float64); ok {
		// sometimes generic json unmarshal makes numbers floats
		val := fmt.Sprintf("%.0f", zpTransIDFloat)
		order.TransactionID = &val
	}

	if errSvc := s.orderRepository.UpdateOrder(ctx, order); errSvc != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "database update failed",
		}, nil
	}

	return &dto.ZaloCallbackResponse{
		ReturnCode:    1,
		ReturnMessage: "success",
	}, nil
}
