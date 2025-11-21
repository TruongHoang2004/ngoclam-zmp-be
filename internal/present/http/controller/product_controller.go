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
	products := r.Group("/products")
	{
		products.POST("", pc.CreateProduct)
		products.GET("/:id", pc.GetProductByID)
		products.GET("", pc.GetAllProduct)
		products.PUT("", pc.UpdateProduct)
		products.DELETE("/:id", pc.DeleteProduct)

		products.POST("/variants", pc.AddProductVariant)
		products.PUT("/variants", pc.UpdateProductVariant)
		products.DELETE("/variants/:id", pc.DeleteProductVariant)
	}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	payload := dto.CreateProductRequest{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.CreateProduct(ctx.Request.Context(), &payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(201, gin.H{"message": "Product created successfully"})
}

func (pc *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	product, err := pc.productService.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, dto.NewProductResponse(product))
}

func (pc *ProductController) GetAllProduct(ctx *gin.Context) {
	request, err := common.ParsePaginationParams(ctx)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid pagination parameters", err))
		return
	}

	products, total, err := pc.productService.ListProducts(ctx.Request.Context(), request.Page, request.Size)
	if err != nil {
		ctx.Error(err)
		return
	}

	productsResponse := make([]*dto.ProductResponse, 0, len(products))
	for _, product := range products {
		productsResponse = append(productsResponse, dto.NewProductResponse(product))
	}

	response := dto.NewPaginationResponse(productsResponse, total, *request)
	ctx.JSON(200, response)
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	payload := dto.UpdateProductRequest{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.UpdateProduct(ctx.Request.Context(), &payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product updated successfully"})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	if err := pc.productService.DeleteProduct(ctx.Request.Context(), id); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product deleted successfully"})
}

func (pc *ProductController) AddProductVariant(ctx *gin.Context) {
	payload := dto.AddProductVariantRequest{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.AddProductVariant(ctx.Request.Context(), &payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(201, gin.H{"message": "Product variant added successfully"})
}

func (pc *ProductController) UpdateProductVariant(ctx *gin.Context) {
	payload := dto.UpdateProductVariantRequest{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.UpdateProductVariant(ctx.Request.Context(), &payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product variant updated successfully"})
}

func (pc *ProductController) DeleteProductVariant(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	if err := pc.productService.DeleteProductVariant(ctx.Request.Context(), id); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product variant deleted successfully"})
}
