package dto

import "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/client/zalo"

type DecodePhoneNumberRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
	Code        string `json:"code" validate:"required"`
}

type DecodePhoneNumberResponse struct {
	PhoneNumber *zalo.UserPhoneNumber `json:"phone_number"`
}

func NewDecodePhoneNumberResponse(phoneNumber *zalo.UserPhoneNumber) *DecodePhoneNumberResponse {
	return &DecodePhoneNumberResponse{
		PhoneNumber: phoneNumber,
	}
}
