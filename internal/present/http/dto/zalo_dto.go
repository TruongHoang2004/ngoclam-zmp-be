package dto

type NofityCallbackData struct {
	OrderID string `json:"orderId"`
	Method  string `json:"method"`
	AppID   string `json:"appId"`
}

type NofityCallbackRequest struct {
	Data NofityCallbackData `json:"data"`
	Mac  string             `json:"mac"`
}

type NofityCallbackResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}

type OrderCallbackRequest struct {
	AppID       string `json:"appId"`
	OrderID     string `json:"orderId"` // Zalo's Order ID
	Method      string `json:"method"`
	Mac         string `json:"mac"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Message     string `json:"message"`
	ResultCode  int    `json:"resultCode"`
	TransID     string `json:"transId"`
	Extradata   string `json:"extradata"`
}

type OrderCallbackResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}
