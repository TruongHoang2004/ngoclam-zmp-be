package services

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/client/zalo"
)

type AuthService struct {
	*baseService
	zaloClient *zalo.Client
	cfg        *config.Config
}

func NewAuthService(
	baseService *baseService,
	zaloClient *zalo.Client,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		baseService: baseService,
		zaloClient:  zaloClient,
		cfg:         cfg,
	}
}

func (s *AuthService) DecodePhoneNumber(ctx context.Context, accessToken string, code string) (*zalo.UserPhoneNumberResponse, *common.Error) {

	phoneNumber, err := s.zaloClient.GetPhoneNumber(ctx, accessToken, code, s.cfg.ZaloAppSecret)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.CurrentService)
	}

	return phoneNumber, nil
}
