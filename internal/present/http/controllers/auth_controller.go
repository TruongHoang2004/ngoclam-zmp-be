package controllers

import (
	"net/http"

	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	*baseController
	authService *services.AuthService
}

func NewAuthController(
	validate *validator.Validate,
	authService *services.AuthService,
) *AuthController {
	return &AuthController{
		baseController: NewBaseController(validate),
		authService:    authService,
	}
}

func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/decode-phone", c.DecodePhoneNumber)
	}
}

func (c *AuthController) DecodePhoneNumber(ctx *gin.Context) {
	var req dto.DecodePhoneNumberRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	phoneNumber, err := c.authService.DecodePhoneNumber(ctx.Request.Context(), req.AccessToken, req.Code)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewDecodePhoneNumberResponse(phoneNumber.Data.Number)))
}
