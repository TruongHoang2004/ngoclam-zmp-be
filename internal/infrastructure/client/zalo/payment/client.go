package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
)

const (
	defaultTimeout = 10 * time.Second
	zaloPaymentURL = "https://payment-mini.zalo.me/api"
)

// ZaloPaymentClient is the Zalo API client.
type ZaloPaymentClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Zalo client.
func NewClient(httpClient *http.Client) *ZaloPaymentClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}
	return &ZaloPaymentClient{
		httpClient: httpClient,
		baseURL:    zaloPaymentURL,
	}
}

// GetOrderStatus
func (c *ZaloPaymentClient) GetOrderStatus(ctx context.Context, config *config.Config, orderId string) (*GetOrderStatusResponse, error) {
	log.Info(ctx, "GetOrderStatus start", "orderID", orderId)

	// Calculate MAC
	// data = "appId={appId}&orderId={orderId}&privateKey={privateKey}"
	dataForMac := fmt.Sprintf("appId=%s&orderId=%s&privateKey=%s",
		config.ZaloAppID, orderId, config.ZaloAppKey)

	mac := utils.ComputeHmac256(dataForMac, config.ZaloAppKey)

	targetURL := c.baseURL + "/transaction/get-status"

	params := url.Values{}
	params.Add("app_id", config.ZaloAppID)
	params.Add("order_id", orderId)
	params.Add("mac", mac)

	targetURLWithParams := targetURL + "?" + params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURLWithParams, nil)
	if err != nil {
		log.Error(ctx, "GetOrderStatus: failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	log.Info(ctx, "GetOrderStatus request", "url", targetURL)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		log.Error(ctx, "GetOrderStatus: failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "GetOrderStatus: unexpected status code", "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response GetOrderStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(ctx, "GetOrderStatus: failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "GetOrderStatus success", "response", response)
	return &response, nil
}

func (c *ZaloPaymentClient) UpdateCodOrderStatus(ctx context.Context, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error) {
	targetURL := fmt.Sprintf("%s/transaction/%s/cod-callback-payment", c.baseURL, req.AppID)
	return c.sendUpdateOrderRequest(ctx, targetURL, req)
}

func (c *ZaloPaymentClient) UpdateBankOrderStatus(ctx context.Context, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error) {
	targetURL := fmt.Sprintf("%s/transaction/%s/bank-callback-payment", c.baseURL, req.AppID)
	return c.sendUpdateOrderRequest(ctx, targetURL, req)
}

func (c *ZaloPaymentClient) sendUpdateOrderRequest(ctx context.Context, url string, req *UpdateOrderStatusRequest) (*UpdateOrderStatusResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Error(ctx, "sendUpdateOrderRequest: failed to marshal request", "error", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(ctx, "sendUpdateOrderRequest: failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	log.Info(ctx, "sendUpdateOrderRequest request", "url", url, "body", string(reqBody))

	resp, err := c.httpClient.Do(request)
	if err != nil {
		log.Error(ctx, "sendUpdateOrderRequest: failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "sendUpdateOrderRequest: unexpected status code", "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response UpdateOrderStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(ctx, "sendUpdateOrderRequest: failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "sendUpdateOrderRequest success", "response", response)
	return &response, nil
}
