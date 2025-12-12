package controllers

import (
	"net/http"

	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	*baseController
	orderService *services.OrderService
}

func NewOrderController(baseController *baseController, orderService *services.OrderService) *OrderController {
	return &OrderController{
		baseController: baseController,
		orderService:   orderService,
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	res, err := c.orderService.CreateOrder(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(res))
}

func (c *OrderController) ListOrders(ctx *gin.Context) {
	pagination, err := c.GetPaginationParams(ctx)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	orders, total, errSvc := c.orderService.ListOrders(ctx.Request.Context(), pagination.Page, pagination.Size)
	if errSvc != nil {
		c.ErrorData(ctx, errSvc)
		return
	}

	// Assuming we return model directly for now or map to DTO if needed.
	// For simplicity, returning model as existing code seems to support JSON tags on models.
	response := dto.NewPaginationResponse(orders, total, *pagination)
	ctx.JSON(http.StatusOK, response)
}

func (c *OrderController) GetOrder(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	order, errSvc := c.orderService.GetOrder(ctx.Request.Context(), id)
	if errSvc != nil {
		c.ErrorData(ctx, errSvc)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(order))
}

func (c *OrderController) RegisterRoutes(r *gin.RouterGroup) {
	orders := r.Group("/orders")
	{
		orders.POST("", c.CreateOrder)
		orders.GET("", c.ListOrders)
		orders.GET("/:id", c.GetOrder)
	}
}
