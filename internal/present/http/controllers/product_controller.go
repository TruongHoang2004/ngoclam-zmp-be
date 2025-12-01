package controllers

import (
	"net/http"

	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
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
	req := dto.CreateProductRequest{}
	if err := pc.BindAndValidateRequest(ctx, &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	if err := pc.productService.CreateProduct(ctx.Request.Context(), &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(nil))
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
	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewProductResponse(product)))
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
	if err := pc.BindAndValidateRequest(ctx, &request); err != nil {
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
	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(nil))
}

func (pc *ProductController) UpdateProductVariant(ctx *gin.Context) {
	req := dto.UpdateProductVariantRequest{}
	if err := pc.BindAndValidateRequest(ctx, &req); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	_, err := pc.productService.UpdateProductVariant(ctx.Request.Context(), &req)
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(nil))
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
	productID, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	images, err := pc.productService.ListProductImages(ctx.Request.Context(), productID)
	if err != nil {
		ctx.Error(err)
		return
	}

	var responses []dto.ProductImageResponse
	for _, image := range images {
		responses = append(responses, *dto.NewProductImageResponse(image))
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(responses))
}

func (pc *ProductController) AttachProductImage(ctx *gin.Context) {
	productID, err := pc.GetUintParam(ctx, "id")
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	payload := dto.AttachProductImageRequest{}
	if err := pc.BindAndValidateRequest(ctx, &payload); err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	img, err := pc.productService.AddProductImage(ctx.Request.Context(), productID, &payload)
	if err != nil {
		pc.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(dto.NewProductImageResponse(img)))
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

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewProductImageResponse(image)))
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
