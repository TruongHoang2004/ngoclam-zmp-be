package services

import (
	"context"
	"encoding/json"
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

func (s *PaymentService) ProcessNotifyCallback(ctx context.Context, req *dto.NofityCallbackRequest) (*dto.NofityCallbackResponse, *common.Error) {
	// 1. Verify Message Authentication Code (HMAC-SHA256)
	// data = 'appId={appId}&orderId={orderId}&method={method}'
	dataForMac := fmt.Sprintf("appId=%s&orderId=%s&method=%s",
		req.Data.AppID, req.Data.OrderID, req.Data.Method)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		// Log for debugging
		fmt.Printf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac)
		return &dto.NofityCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Order ID
	orderIDUint64, err := strconv.ParseUint(req.Data.OrderID, 10, 32)
	if err != nil {
		return &dto.NofityCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid order id format",
		}, nil
	}
	orderID := uint(orderIDUint64)

	// 3. Update Order Satus
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		return &dto.NofityCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "order not found",
		}, nil
	}

	// Check if already paid or final state
	if order.Status == "success" {
		return &dto.NofityCallbackResponse{
			ReturnCode:    1,
			ReturnMessage: "success",
		}, nil
	}

	// Update status
	order.Status = "success"
	order.TransactionID = &req.Data.OrderID

	if errSvc := s.orderRepository.UpdateOrder(ctx, order); errSvc != nil {
		return &dto.NofityCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "database update failed",
		}, nil
	}

	return &dto.NofityCallbackResponse{
		ReturnCode:    1,
		ReturnMessage: "success",
	}, nil
}

func (s *PaymentService) ProcessOrderCallback(ctx context.Context, req *dto.OrderCallbackRequest) (*dto.OrderCallbackResponse, *common.Error) {
	// 1. Verify Message Authentication Code (HMAC-SHA256)
	// dataForMac = "appId={appId}&amount={amount}&description={description}&orderId={orderId}&message={message}&resultCode={resultCode}&transId={transId}"
	dataForMac := fmt.Sprintf("appId=%s&amount=%d&description=%s&orderId=%s&message=%s&resultCode=%d&transId=%s",
		req.AppID, req.Amount, req.Description, req.OrderID, req.Message, req.ResultCode, req.TransID)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		fmt.Printf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac)
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Extradata to get internal Order ID
	var extraDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(req.Extradata), &extraDataMap); err != nil {
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid extradata format",
		}, nil
	}

	pkOrderIDFloat, ok := extraDataMap["pk_order_id"].(float64)
	if !ok {
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "pk_order_id not found in extradata",
		}, nil
	}
	orderID := uint(pkOrderIDFloat)

	// 3. Update Order Status
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "order not found",
		}, nil
	}

	// Check if already paid
	if order.Status == "success" {
		// Idempotency: return success if already success
		return &dto.OrderCallbackResponse{
			ReturnCode:    1,
			ReturnMessage: "success",
		}, nil
	}

	if req.ResultCode == 1 {
		order.Status = "success"
	} else {
		// Only update to failed if currently pending? or just log?
		// For now mark as failed
		order.Status = "failed"
	}

	// Store Zalo Trans ID if not already
	order.TransactionID = &req.TransID

	if errSvc := s.orderRepository.UpdateOrder(ctx, order); errSvc != nil {
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "database update failed",
		}, nil
	}

	return &dto.OrderCallbackResponse{
		ReturnCode:    1,
		ReturnMessage: "success",
	}, nil
}
