package controllers

import (
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
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

func (c *PaymentController) NotifyCallback(ctx *gin.Context) {

	var req dto.NofityCallbackRequest

	log.Info(ctx.Request.Context(), "NotifyCallback: %v", req)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx.Request.Context()).SetDetail(err.Error()))
		return
	}

	res, err := c.paymentService.ProcessNotifyCallback(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *PaymentController) OrderCallback(ctx *gin.Context) {

	var req dto.OrderCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx.Request.Context()).SetDetail(err.Error()))
		return
	}

	log.Info(ctx.Request.Context(), "OrderCallback: %v", req)

	res, err := c.paymentService.ProcessOrderCallback(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *PaymentController) RegisterRoutes(r *gin.RouterGroup) {
	payment := r.Group("/payment")
	{
		payment.POST("/notify-callback", c.NotifyCallback)
		payment.POST("/order-callback", c.OrderCallback)
	}
}
