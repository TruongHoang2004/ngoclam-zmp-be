package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
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
	// dataForMac = "appId={appId}&amount={amount}&description={description}&orderId={orderId}&message={message}&resultCode={resultCode}&transId={transId}"
	dataForMac := fmt.Sprintf("appId=%s&amount=%d&description=%s&orderId=%s&message=%s&resultCode=%d&transId=%s",
		req.AppID, req.Amount, req.Description, req.OrderID, req.Message, req.ResultCode, req.TransID)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		// Log for debugging?
		fmt.Printf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac)
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Parse Extradata to get Order ID
	// Extradata: {"pk_order_id": 123}
	var extraDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(req.Extradata), &extraDataMap); err != nil {
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "invalid extradata format",
		}, nil
	}

	pkOrderIDFloat, ok := extraDataMap["pk_order_id"].(float64)
	if !ok {
		// Fallback to string check?
		return &dto.ZaloCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "pk_order_id not found in extradata",
		}, nil
	}
	orderID := uint(pkOrderIDFloat)

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

	if req.ResultCode == 1 {
		order.Status = "success"
	} else {
		order.Status = "failed" // or other status
	}

	// Store Zalo Trans ID
	order.TransactionID = &req.TransID

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
