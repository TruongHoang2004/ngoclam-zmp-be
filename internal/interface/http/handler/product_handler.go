package handler

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/application"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/interface/http/dto"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService *application.ProductService
}

func NewProductHandler(productService *application.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("/:id", h.GetProductByID)
		products.GET("", h.GetAllProducts)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("/category/:category_id", h.GetProductsByCategoryID)
	}
}

// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product to create"
// @Success 200 {object} dto.ProductResponseDTO "Returns the created product"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var createProductRequest dto.CreateProductRequest

	// Validate the request
	if err := ctx.ShouldBindJSON(&createProductRequest); err != nil {
		application.HandleError(ctx, err)
		return
	}

	product, err := h.ProductService.CreateProduct(ctx, *createProductRequest.ToDomain())
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, dto.NewProductResponseDTO(*product))
}

// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponseDTO "Returns the product"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	id, err := ConvertStringToUint(ctx.Param("id"))
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	product, err := h.ProductService.GetProductByID(ctx, id)
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Product retrieved successfully", "product": dto.NewProductResponseDTO(*product)})
}

// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} dto.ProductResponseDTO "Returns the list of products"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := h.ProductService.GetAllProducts(ctx)
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	var productDTOs []dto.ProductResponseDTO
	for _, p := range products {
		productDTOs = append(productDTOs, dto.NewProductResponseDTO(*p))
	}

	ctx.JSON(200, gin.H{"message": "Products retrieved successfully", "products": productDTOs})
}

// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product to update"
// @Success 200 {object} dto.ProductResponseDTO "Returns the updated product"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {

	var updateProductRequest dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&updateProductRequest); err != nil {
		application.HandleError(ctx, err)
		return
	}

	product, err := h.ProductService.UpdateProduct(ctx.Request.Context(), *updateProductRequest.ToDomain())
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Product updated successfully", "product": dto.NewProductResponseDTO(*product)})
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "Returns a success message"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id, err := ConvertStringToUint(ctx.Param("id"))
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	err = h.ProductService.DeleteProduct(ctx.Request.Context(), id)
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	ctx.JSON(200, gin.H{"message": "Product deleted successfully"})
}

// @Summary Get products by category ID
// @Description Get a list of products by category ID
// @Tags products
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {array} dto.ProductResponseDTO "Returns the list of products in the category"
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetProductsByCategoryID(ctx *gin.Context) {
	categoryID, err := ConvertStringToUint(ctx.Param("category_id"))
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	products, err := h.ProductService.GetProductsByCategoryID(ctx.Request.Context(), categoryID)
	if err != nil {
		application.HandleError(ctx, err)
		return
	}

	var productDTOs []dto.ProductResponseDTO
	for _, p := range products {
		productDTOs = append(productDTOs, dto.NewProductResponseDTO(*p))
	}

	ctx.JSON(200, gin.H{"message": "Products retrieved successfully", "data": productDTOs})
}
