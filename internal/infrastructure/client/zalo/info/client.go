package info

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
)

const (
	defaultTimeout = 10 * time.Second
	zaloGraphURL   = "https://graph.zalo.me/v2.0/me/info"
)

// ZaloInfoClient is the Zalo API client.
type ZaloInfoClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Zalo client.
func NewClient(httpClient *http.Client) *ZaloInfoClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}
	return &ZaloInfoClient{
		httpClient: httpClient,
		baseURL:    zaloGraphURL,
	}
}

// GetPhoneNumber retrieves the user's phone number from Zalo API using the provided token.
func (c *ZaloInfoClient) GetPhoneNumber(ctx context.Context, accessToken, code, secretKey string) (*UserPhoneNumberResponse, error) {
	log.Info(ctx, "GetPhoneNumber start", "baseURL", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		log.Error(ctx, "GetPhoneNumber: failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("access_token", accessToken)
	req.Header.Set("code", code)
	req.Header.Set("secret_key", secretKey)

	log.Info(ctx, "GetPhoneNumber request", "request", req)
	log.Info(ctx, "GetPhoneNumber request header", "header", req.Header)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(ctx, "GetPhoneNumber: failed to send request", "error", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "GetPhoneNumber: unexpected status code", "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var phoneNumber UserPhoneNumberResponse
	if err := json.NewDecoder(resp.Body).Decode(&phoneNumber); err != nil {
		log.Error(ctx, "GetPhoneNumber: failed to decode response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "GetPhoneNumber success", "phoneNumber", phoneNumber)
	return &phoneNumber, nil
}
