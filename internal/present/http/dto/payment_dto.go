package dto

type WebhookReceiverRequest struct {
	ID              int64   `json:"id"`
	Gateway         string  `json:"gateway"`
	TransactionDate string  `json:"transactionDate"`
	AccountNumber   string  `json:"accountNumber"`
	Code            *string `json:"code"`
	Content         string  `json:"content"`
	TransferType    string  `json:"transferType"`
	TransferAmount  int64   `json:"transferAmount"`
	Accumulated     int64   `json:"accumulated"`
	SubAccount      *string `json:"subAccount"`
	ReferenceCode   string  `json:"referenceCode"`
	Description     string  `json:"description"`
}
