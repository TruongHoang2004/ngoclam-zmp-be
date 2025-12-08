package zalo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
	zaloGraphURL   = "https://graph.zalo.me/v2.0/me/info"
)

// Client is the Zalo API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Zalo client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    zaloGraphURL,
	}
}

// GetUserInfo retrieves the user information from Zalo API.
func (c *Client) GetUserInfo(userAccessToken, token, secretKey string) (*UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("access_token", userAccessToken)
	req.Header.Set("code", token)
	req.Header.Set("secret_key", secretKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &userInfo, nil
}
