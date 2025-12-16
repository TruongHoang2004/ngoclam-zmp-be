package payment

// UpdateOrderStatusRequest represents the request payload for updating order status.
type UpdateOrderStatusRequest struct {
	AppID      string `json:"appId"`
	OrderID    string `json:"orderId"`
	ResultCode int    `json:"resultCode"` // 1: Success, 0: Refunded, -1: Failed
	Mac        string `json:"mac"`
}

// UpdateOrderStatusResponse represents the response from Zalo API.
type UpdateOrderStatusResponse struct {
	Error int `json:"error"`
	Data  struct {
		ReturnCode    int    `json:"returnCode"`
		ReturnMessage string `json:"returnMessage"`
	} `json:"data"`
}

type GetOrderStatusRequest struct {
	AppID   string `json:"appId"`
	OrderID string `json:"orderId"`
	Mac     string `json:"mac"`
}

type GetOrderStatusResponse struct {
	Error int `json:"error"`
	Data  struct {
		ReturnCode      int    `json:"returnCode"`
		ReturnMessage   string `json:"returnMessage"`
		IsProcessing    bool   `json:"isProcessing"`
		TransID         string `json:"transId"`
		Method          string `json:"method"`
		Amount          int64  `json:"amount"`
		TransTime       int64  `json:"transTime"`
		MerchantTransID string `json:"merchantTransId"`
		ExtraData       string `json:"extraData"`
	} `json:"data"`
}
