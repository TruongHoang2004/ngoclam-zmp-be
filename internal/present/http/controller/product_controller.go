package controller

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (pc *ProductController) RegisterRoutes(r *gin.RouterGroup) {

}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	dto := dto.CreateProductRequest{}
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	pc.productService.CreateProduct(ctx.Request.Context(), nil)
}

func (pc *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	product, err := pc.productService.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(common.NotFound("Product not found", err))
		return
	}
	ctx.JSON(200, dto.NewProductResponse(product))
}
