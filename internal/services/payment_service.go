package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type PaymentMethod string

const (
	PaymentMethodBank PaymentMethod = "BANK"
	PaymentMethodCod  PaymentMethod = "COD"
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
	// data = 'appId={appId}&orderId={orderId}&method={method}'
	dataForMac := fmt.Sprintf("appId=%s&orderId=%s&method=%s",
		req.Data.AppID, req.Data.OrderID, req.Data.Method)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		// Log for debugging
		fmt.Printf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac)
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Order ID
	orderIDUint64, err := strconv.ParseUint(req.Data.OrderID, 10, 32)
	if err != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid order id format",
		}, nil
	}
	orderID := uint(orderIDUint64)

	// 3. Update Order Satus
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "order not found",
		}, nil
	}

	// Check if already paid or final state
	if order.Status == "success" {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    1,
			ReturnMessage: "success",
		}, nil
	}

	// Update status
	order.Status = "success"

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
