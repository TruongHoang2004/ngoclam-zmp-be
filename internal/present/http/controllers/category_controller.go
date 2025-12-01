package controllers

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryController struct {
	*baseController
	categoryService *services.CategoryService
}

func NewCategoryController(
	validate *validator.Validate,
	categoryService *services.CategoryService,
) *CategoryController {
	return &CategoryController{
		baseController:  NewBaseController(validate),
		categoryService: categoryService,
	}
}

func (c *CategoryController) RegisterRoutes(r *gin.RouterGroup) {
	categories := r.Group("/categories")
	{
		categories.POST("", c.CreateCategory)
		categories.GET("/:id", c.GetCategoryByID)
		categories.GET("", c.ListCategories)
		categories.PUT("/:id", c.UpdateCategory)
		categories.DELETE("/:id", c.DeleteCategory)
		categories.GET("/:id/products", c.GetProductsByCategory)
	}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	res, err := c.categoryService.CreateCategory(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, res)
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	res, err := c.categoryService.GetCategoryByID(ctx.Request.Context(), id)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, res)
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
	res, err := c.categoryService.ListCategories(ctx.Request.Context())
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, res)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}
	req.ID = id

	res, err := c.categoryService.UpdateCategory(ctx.Request.Context(), &req)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, res)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	if err := c.categoryService.DeleteCategory(ctx.Request.Context(), id); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, map[string]string{"message": "success"})
}

func (c *CategoryController) GetProductsByCategory(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	paginationReq, err := c.GetPaginationParams(ctx)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	res, err := c.categoryService.GetProductsByCategory(ctx.Request.Context(), id, paginationReq.Page, paginationReq.Size)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	c.Success(ctx, res)
}
