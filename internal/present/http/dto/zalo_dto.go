package dto

type ZaloCallbackData struct {
	AppID   string `json:"appId"`
	OrderID string `json:"orderId"` // Zalo's Order ID
	Method  string `json:"method"`
	Mac     string `json:"mac"`
}

type ZaloCallbackRequest struct {
	Data ZaloCallbackData `json:"data"`
	Mac  string           `json:"mac"`
}

type ZaloCallbackResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}
