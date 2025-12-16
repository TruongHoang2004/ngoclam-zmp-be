package controllers

import (
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	*baseController
	paymentService *services.PaymentService
}

func NewPaymentController(baseController *baseController, paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		baseController: baseController,
		paymentService: paymentService,
	}
}

func (c *PaymentController) ZaloCallback(ctx *gin.Context) {

	var req dto.ZaloCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx.Request.Context()).SetDetail(err.Error()))
		return
	}

	res, err := c.paymentService.ProcessZaloCallback(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *PaymentController) RegisterRoutes(r *gin.RouterGroup) {
	payment := r.Group("/payment")
	{
		payment.POST("/zalo-callback", c.ZaloCallback)
	}
}
