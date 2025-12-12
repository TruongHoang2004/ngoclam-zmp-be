package dto

type ZaloCallbackRequest struct {
	AppID           string `json:"appId"`
	OrderID         string `json:"orderId"` // Zalo's Order ID
	TransID         string `json:"transId"`
	Method          string `json:"method"`
	TransTime       int64  `json:"transTime"`
	MerchantTransID string `json:"merchantTransId"`
	Amount          int64  `json:"amount"`
	Description     string `json:"description"`
	ResultCode      int    `json:"resultCode"`
	Message         string `json:"message"`
	Extradata       string `json:"extradata"`
	Mac             string `json:"mac"`
}

type ZaloCallbackResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}
