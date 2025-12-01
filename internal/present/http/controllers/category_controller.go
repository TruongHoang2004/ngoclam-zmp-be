package controllers

import (
	"net/http"
	"strconv"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
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
		categories.GET("/:id", c.GetCategory)
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

	if err := c.categoryService.CreateCategory(ctx, &req); err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(nil))
}

func (c *CategoryController) GetCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrBadRequest(ctx).SetDetail("Invalid category ID"))
		return
	}

	category, serviceErr := c.categoryService.GetCategoryByID(ctx, uint(id))
	if serviceErr != nil {
		ctx.JSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewCategoryResponse(category)))
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
	categories, err := c.categoryService.ListCategories(ctx)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	var responses []dto.CategoryResponse
	for _, cat := range categories {
		responses = append(responses, *dto.NewCategoryResponse(cat))
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(responses))
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

	updatedCategory, err := c.categoryService.UpdateCategory(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewCategoryResponse(updatedCategory)))
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
