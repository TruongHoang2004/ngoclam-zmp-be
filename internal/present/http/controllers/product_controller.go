package controllers

import (
	"net/http"

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
	id, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	product, err := pc.productService.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(200, dto.NewProductResponse(product))
}

func (pc *ProductController) GetAllProduct(ctx *gin.Context) {
	pagination, err := pc.GetPaginationParams(ctx)
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	products, total, err := pc.productService.ListProducts(ctx.Request.Context(), pagination.Page, pagination.Size)
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	productsResponse := make([]*dto.ProductResponse, 0, len(products))
	for _, product := range products {
		productsResponse = append(productsResponse, dto.NewProductResponse(product))
	}

	response := dto.NewPaginationResponse(productsResponse, total, *pagination)
	ctx.JSON(200, response)
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	request := dto.UpdateProductRequest{}
	if err := pc.BindAndValidateRequest(ctx, request); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.UpdateProduct(ctx.Request.Context(), &request); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product updated successfully"})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.DeleteProduct(ctx.Request.Context(), id); err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product deleted successfully"})
}

func (pc *ProductController) AddProductVariant(ctx *gin.Context) {
	req := dto.AddProductVariantRequest{}
	if err := pc.BindAndValidateRequest(ctx, &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.AddProductVariant(ctx.Request.Context(), &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product variant added successfully"})
}

func (pc *ProductController) UpdateProductVariant(ctx *gin.Context) {
	req := dto.UpdateProductVariantRequest{}
	if err := pc.BindAndValidateRequest(ctx, &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.UpdateProductVariant(ctx.Request.Context(), &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product variant updated successfully"})
}

func (pc *ProductController) DeleteProductVariant(ctx *gin.Context) {
	id, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.DeleteProductVariant(ctx.Request.Context(), id); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.DeleteProductVariant(ctx.Request.Context(), id); err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Product variant deleted successfully"})
}

func (pc *ProductController) ListProductImages(ctx *gin.Context) {
	id, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
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
	id, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	payload := dto.AttachProductImageRequest{}
	if err := pc.BindAndValidateRequest(ctx, &payload); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	image, err := pc.productService.AddProductImage(ctx.Request.Context(), id, &payload)
	if err != nil {
		pc.ErrorData(ctx, err)
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
		pc.ErrorData(ctx, err)
		return
	}

	ctx.JSON(200, dto.NewProductImageResponse(image))
}

func (pc *ProductController) DeleteProductImage(ctx *gin.Context) {
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

	if err := pc.productService.DeleteProductImage(ctx.Request.Context(), productID, imageID); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Product image deleted successfully"})
}
