package controllers

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	*baseController
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
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

		products.GET("/:id/images", pc.ListProductImages)
		products.POST("/:id/images", pc.AttachProductImage)
		products.PATCH("/:id/images/:imageId", pc.UpdateProductImage)
		products.DELETE("/:id/images/:imageId", pc.DeleteProductImage)
	}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	req := new(dto.CreateProductRequest)
	if err := pc.BindAndValidateRequest(ctx, req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.CreateProduct(ctx.Request.Context(), req); err != nil {
		pc.ErrorData(ctx, err)
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.UpdateProduct(ctx.Request.Context(), &req); err != nil {
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.AddProductVariant(ctx.Request.Context(), &req); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(201, gin.H{"message": "Product variant added successfully"})
}

func (pc *ProductController) UpdateProductVariant(ctx *gin.Context) {
	payload := dto.UpdateProductVariantRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	if err := pc.productService.UpdateProductVariant(ctx.Request.Context(), &req); err != nil {
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

func (pc *ProductController) ListProductImages(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	images, err := pc.productService.ListProductImages(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	imageResponses := make([]*dto.ProductImageResponse, 0, len(images))
	for _, img := range images {
		imageResponses = append(imageResponses, dto.NewProductImageResponse(img))
	}

	ctx.JSON(200, gin.H{"images": imageResponses})
}

func (pc *ProductController) AttachProductImage(ctx *gin.Context) {
	id, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	payload := dto.AttachProductImageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	image, err := pc.productService.AddProductImage(ctx.Request.Context(), id, &payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, dto.NewProductImageResponse(image))
}

func (pc *ProductController) UpdateProductImage(ctx *gin.Context) {
	productID, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	imageID, err := pc.GetUintParam(ctx, "imageId")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	req := dto.UpdateProductImageRequest{}
	if err := pc.BindAndValidateRequest(ctx, &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	image, err := pc.productService.UpdateProductImage(ctx.Request.Context(), productID, imageID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, dto.NewProductImageResponse(image))
}

func (pc *ProductController) DeleteProductImage(ctx *gin.Context) {
	productID, err := common.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	imageID, err := common.ParseUintParam(ctx, "imageId")
	if err != nil {
		ctx.Error(common.BadRequest("Invalid image ID format", err))
		return
	}

	if err := pc.productService.DeleteProductImage(ctx.Request.Context(), productID, imageID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Product image deleted successfully"})
}
