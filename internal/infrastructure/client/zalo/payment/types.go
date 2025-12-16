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
	Msg  string `json:"msg"`
	Err  int    `json:"err"`
	Data struct {
		ReturnCode      int    `json:"returnCode"`
		ReturnMessage   string `json:"returnMessage"`
		IsProcessing    bool   `json:"isProcessing"`
		TransID         string `json:"transId"`
		Method          string `json:"method"`
		Amount          int64  `json:"amount"`
		TransTime       int64  `json:"transTime"`
		MerchantTransID int64  `json:"merchantTransId"`
		Extradata       string `json:"extradata"`
		SubResultCode   int    `json:"subResultCode"`
		UpdateAt        int64  `json:"updateAt"`
		CreateAt        int64  `json:"createAt"`
	} `json:"data"`
}
