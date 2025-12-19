package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/client/zalo/payment"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/model"
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

	return &dto.NofityCallbackResponse{
		ReturnCode:    1,
		ReturnMessage: "success",
	}, nil
}

func (s *PaymentService) deferredCheckOrderStatus(zaloOrderID string) {
	// Wait for 5 minutes
	time.Sleep(5 * time.Minute)

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
	if orderStatus.Err != 0 {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: order %s error. Error: %d\n",
			zaloOrderID, orderStatus.Err))
		return
	}

	currentStatus := model.OrderStatusPending
	if orderStatus.Data.ReturnCode == 1 {
		currentStatus = model.OrderStatusCompleted
	}
	if orderStatus.Data.ReturnCode == 0 {
		currentStatus = model.OrderStatusProcessing
	}
	if orderStatus.Data.ReturnCode == -1 {
		currentStatus = model.OrderStatusFailed
	}

	// Parse ExtraData to get internal Order ID
	decodedExtradata, err := url.QueryUnescape(orderStatus.Data.Extradata)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to unescape extradata for %s: %v\n", zaloOrderID, err))
		return
	}

	var extraDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(decodedExtradata), &extraDataMap); err != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to parse extradata for %s: %v\n", zaloOrderID, err))
		return
	}

	pkOrderID, ok := extraDataMap["pk_order_id"].(string)
	if !ok {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: pk_order_id not found in extradata for %s\n", zaloOrderID))
		return
	}
	orderID := pkOrderID

	// Update Order Status in DB
	order, errSvc := s.orderRepository.GetOrder(ctx, orderID)
	if errSvc != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to get order %s: %v\n", orderID, errSvc))
		return
	}

	// Idempotency check
	if order.Status == string(model.OrderStatusCompleted) {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: order %s already success\n", orderID))
		return
	}

	order.Status = string(currentStatus)
	// Ensure TransactionID is set if missing
	if order.TransactionID == nil || *order.TransactionID == "" {
		transID := orderStatus.Data.TransID // Assuming TransID is available in status response
		order.TransactionID = &transID
	}

	if errSvc := s.orderRepository.UpdateOrder(ctx, order); errSvc != nil {
		log.Error(ctx, fmt.Sprintf("deferredCheckOrderStatus: failed to update order %s: %v\n", orderID, errSvc))
		return
	}

	log.Info(ctx, fmt.Sprintf("deferredCheckOrderStatus: successfully updated order %s to success\n", orderID))
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

	pkOrderID, ok := extraDataMap["pk_order_id"].(string)
	if !ok {
		return &dto.OrderCallbackResponse{
			ReturnCode:    -1,
			ReturnMessage: "pk_order_id not found in extradata",
		}, nil
	}
	orderID := pkOrderID

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

func (s *PaymentService) ProcessWebhookReceiver(ctx context.Context, req *dto.WebhookReceiverRequest) *common.Error {

	log.Debug(ctx, fmt.Sprintf("ProcessWebhookReceiver: received content: %s", req.Content))

	contentParts := strings.Fields(req.Content)
	if len(contentParts) == 0 {
		return common.ErrBadRequest(ctx)
	}
	responseOrderID := contentParts[0]
	log.Debug(ctx, fmt.Sprintf("ProcessWebhookReceiver: parsed order ID: %s", responseOrderID))

	order, errSvc := s.orderRepository.GetOrder(ctx, responseOrderID)
	if errSvc != nil {
		return common.ErrNotFound(ctx, "Order", "not found")
	}

	order.Status = string(model.OrderStatusPaid)

	s.orderRepository.UpdateOrder(ctx, order)

	// Notify Zalo Mini App
	resultCode := 1
	dataForMac := fmt.Sprintf("appId=%s&orderId=%s&resultCode=%d", s.cfg.ZaloAppID, order.ID, resultCode)
	mac := utils.ComputeHmac256(dataForMac, s.cfg.ZaloAppSecret)

	payload, _ := json.Marshal(map[string]interface{}{
		"appId":      s.cfg.ZaloAppID,
		"orderId":    order.ZaloOrderID,
		"resultCode": resultCode,
		"mac":        mac,
	})

	zaloURL := fmt.Sprintf("https://payment-mini.zalo.me/api/transaction/%s/bank-callback-payment", s.cfg.ZaloAppID)
	log.Debug(ctx, fmt.Sprintf("ProcessWebhookReceiver: sending to Zalo Mini App URL: %s, payload: %s", zaloURL, string(payload)))

	resp, err := http.Post(zaloURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return common.ErrSystemError(ctx, "Notify Zalo Mini App failed")
	}

	defer resp.Body.Close()

	log.Debug(ctx, fmt.Sprintf("Notify Zalo Mini App response status: %s", resp.Status))
	log.Debug(ctx, fmt.Sprintf("Notify Zalo Mini App success: %s\n", order.ID))

	// 2. Check order status after 5 minutes
	s.deferredCheckOrderStatus(*order.ZaloOrderID)

	return nil
}
