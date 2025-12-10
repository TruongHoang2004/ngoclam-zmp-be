package dto

type DecodePhoneNumberRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
	Code        string `json:"code" validate:"required"`
}

type DecodePhoneNumberResponse struct {
	PhoneNumber string `json:"phone_number"`
}

func NewDecodePhoneNumberResponse(phoneNumber string) *DecodePhoneNumberResponse {
	return &DecodePhoneNumberResponse{
		PhoneNumber: phoneNumber,
	}
}
