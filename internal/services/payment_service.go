package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/client/zalo/payment"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
)

type PaymentMethod string

const (
	PaymentMethodBank PaymentMethod = "BANK"
	PaymentMethodCod  PaymentMethod = "COD"
)

type PaymentService struct {
	orderRepository   *repositories.OrderRepository
	zaloPaymentClient *payment.ZaloPaymentClient
	cfg               *config.Config
}

func NewPaymentService(orderRepo *repositories.OrderRepository, zaloPaymentClient *payment.ZaloPaymentClient, cfg *config.Config) *PaymentService {
	return &PaymentService{
		orderRepository:   orderRepo,
		zaloPaymentClient: zaloPaymentClient,
		cfg:               cfg,
	}
}

func (s *PaymentService) ProcessNotifyCallback(ctx context.Context, req *dto.NofityCallbackRequest) (*dto.NofityCallbackResponse, *common.Error) {

	// 1. Verify Message Authentication Code (HMAC-SHA256)
	// data = 'appId={appId}&orderId={orderId}&method={method}'
	dataForMac := fmt.Sprintf("appId=%s&orderId=%s&method=%s",
		req.Data.AppID, req.Data.OrderID, req.Data.Method)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppPrivateKey)
	if mac != req.Mac {
		// Log for debugging
		log.Debug(ctx, fmt.Sprintf("dataForMac: %s\n", dataForMac))
		log.Error(ctx, fmt.Sprintf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac))
		return &dto.NofityCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "mac not equal",
		}, nil
	}

	// 2. Start background job to check status after 5 minutes
	go s.deferredCheckOrderStatus(req.Data.OrderID)

	return &dto.NofityCallbackResponse{
		ReturnCode:    1,
		ReturnMessage: "success",
	}, nil
}

func (s *PaymentService) deferredCheckOrderStatus(zaloOrderID string) {
	// Wait for 1 minutes
	time.Sleep(1 * time.Minute)

	ctx := context.Background()
	log.Debug(ctx, fmt.Sprintf("Starting deferred check for Zalo Order ID: %s\n", zaloOrderID))

	// Get Order Status from Zalo
	orderStatus, err := s.zaloPaymentClient.GetOrderStatus(ctx, s.cfg, zaloOrderID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to get status for %s: %v\n", zaloOrderID, err))
		return
	}

	log.Debug(ctx, fmt.Sprintf("deferredCheckOrderStatus: order %s status: %v\n", zaloOrderID, orderStatus))

	// Only proceed if payment is successful
	if orderStatus.Error != 0 || orderStatus.Data.ReturnCode != 1 {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: order %s not successful or error. ReturnCode: %d, Error: %d\n",
			zaloOrderID, orderStatus.Data.ReturnCode, orderStatus.Error))
		return
	}

	// Parse ExtraData to get internal Order ID
	var extraDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(orderStatus.Data.ExtraData), &extraDataMap); err != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to parse extradata for %s: %v\n", zaloOrderID, err))
		return
	}

	pkOrderIDFloat, ok := extraDataMap["pk_order_id"].(float64)
	if !ok {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: pk_order_id not found in extradata for %s\n", zaloOrderID))
		return
	}
	orderID := uint(pkOrderIDFloat)

	// Update Order Status in DB
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to get order %d: %v\n", orderID, errSvc))
		return
	}

	// Idempotency check
	if order.Status == "success" {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: order %d already success\n", orderID))
		return
	}

	order.Status = "success"
	// Ensure TransactionID is set if missing
	if order.TransactionID == nil || *order.TransactionID == "" {
		transID := orderStatus.Data.TransID // Assuming TransID is available in status response
		order.TransactionID = &transID
	}

	if errSvc := s.orderRepository.UpdateOrder(ctx, order); errSvc != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to update order %d: %v\n", orderID, errSvc))
		return
	}

	log.Info(ctx, fmt.Sprintf("deferredCheckOrderStatus: successfully updated order %d to success\n", orderID))
}

func (s *PaymentService) ProcessOrderCallback(ctx context.Context, req *dto.OrderCallbackRequest) (*dto.OrderCallbackResponse, *common.Error) {
	// 1. Verify Message Authentication Code (HMAC-SHA256)
	// dataForMac = "appId={appId}&amount={amount}&description={description}&orderId={orderId}&message={message}&resultCode={resultCode}&transId={transId}"
	dataForMac := fmt.Sprintf("appId=%s&amount=%d&description=%s&orderId=%s&message=%s&resultCode=%d&transId=%s",
		req.AppID, req.Amount, req.Description, req.OrderID, req.Message, req.ResultCode, req.TransID)

	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)
	if mac != req.Mac {
		log.Error(ctx, fmt.Sprintf("Invalid MAC: calculated %s, received %s\n", mac, req.Mac))
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
