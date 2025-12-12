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
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type PaymentService struct {
	orderRepository *repositories.OrderRepository
	cfg             *config.Config
}

func NewPaymentService(orderRepo *repositories.OrderRepository, cfg *config.Config) *PaymentService {
	return &PaymentService{
		orderRepository: orderRepo,
		cfg:             cfg,
	}
}

func (s *PaymentService) ProcessZaloCallback(ctx context.Context, req *dto.ZaloCallbackRequest) (*dto.ZaloCallbackResponse, *common.Error) {
	// 1. Verify Message Authentication Code (HMAC-SHA256)
	mac := utils.ComputeHmac256(req.Data, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Data
	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(req.Data), &dataMap); err != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid data format",
		}, nil
	}

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
