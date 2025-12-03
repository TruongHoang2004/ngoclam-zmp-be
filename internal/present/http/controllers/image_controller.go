package controllers

import (
	"net/http"

	"io"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	*baseController
	imageService *services.ImageService
}

func NewImageController(base *baseController, imageService *services.ImageService) *ImageController {
	return &ImageController{
		baseController: base,
		imageService:   imageService,
	}
}

func (c *ImageController) RegisterRoutes(r *gin.RouterGroup) {
	images := r.Group("/images")
	{
		images.POST("", c.UploadImage)
		images.POST("/url", c.UploadImageFromURL)
		images.GET("/:id", c.GetImageByID)
		images.GET("", c.GetAllImages)
		images.PUT("/:id", c.UpdateImage)
		images.PUT("/:id/url", c.UpdateImageFromURL)
		images.DELETE("/:id", c.DeleteImage)
	}
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	fileHeader, err := c.GetFile(ctx, "file")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	f, openErr := fileHeader.Open()
	if openErr != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetDetail(openErr.Error()).SetSource(common.CurrentService))
		return
	}
	defer f.Close()

	fileBytes, readErr := io.ReadAll(f)
	if readErr != nil {
		c.ErrorData(ctx, common.ErrSystemError(ctx, readErr.Error()))
		return
	}

	// The original filename uploaded by the client
	fileName := fileHeader.Filename

	// Use io.Reader directly for better memory efficiency
	image, err := c.imageService.UploadImage(ctx, fileName, fileBytes)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(dto.NewImageResponse(image)))
}

// UploadImageFromURL handles image upload from a URL
func (c *ImageController) UploadImageFromURL(ctx *gin.Context) {
	var req struct {
		URL      string `json:"url" validator:"required"`
		FileName string `json:"file_name,omitempty"`
	}

	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	image, err := c.imageService.UploadImageFromURL(ctx, req.URL, req.FileName)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(dto.NewImageResponse(image)))
}

func (c *ImageController) GetImageByID(ctx *gin.Context) {

	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	image, err := c.imageService.GetImageByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewImageResponse(image)))
}

func (c *ImageController) GetAllImages(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetSource(common.CurrentService))
		return
	}
	req.Normalize()

	imageList, total, err := c.imageService.GetAllImages(ctx, req.Page, req.Size)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	var responses []dto.ImageResponse
	for _, img := range imageList {
		responses = append(responses, *dto.NewImageResponse(img))
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewPaginationResponse(responses, total, dto.PaginationRequest{Page: req.Page, Size: req.Size})))

}

func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	fileHeader, err := c.GetFile(ctx, "file")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	f, openErr := fileHeader.Open()
	if openErr != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetDetail(openErr.Error()).SetSource(common.CurrentService))
		return
	}
	defer f.Close()

	fileBytes, readErr := io.ReadAll(f)
	if readErr != nil {
		c.ErrorData(ctx, common.ErrSystemError(ctx, readErr.Error()))
		return
	}

	fileName := fileHeader.Filename

	image, err := c.imageService.UpdateImage(ctx, uint(id), fileName, fileBytes)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewImageResponse(image)))
}

// UpdateImageFromURL handles image update from a URL
func (c *ImageController) UpdateImageFromURL(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	var req struct {
		URL      string `json:"url" binding:"required"`
		FileName string `json:"file_name,omitempty"`
	}

	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	image, err := c.imageService.UpdateImageFromURL(ctx, uint(id), req.URL, req.FileName)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewImageResponse(image)))
}

func (c *ImageController) DeleteImage(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	err = c.imageService.DeleteImage(ctx.Request.Context(), id)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
