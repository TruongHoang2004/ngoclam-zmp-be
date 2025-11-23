package controllers

import (
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
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
	file, err := c.GetFile(ctx, "file")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	// The original filename uploaded by the client
	fileName := file.Filename

	// Open the file to read its content
	f, ferr := file.Open()
	if ferr != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetDetail(ferr.Error()).SetSource(common.CurrentService))
		return
	}
	defer f.Close()

	// Use io.Reader directly for better memory efficiency
	image, err := c.imageService.UploadImageFromReader(ctx.Request.Context(), f, fileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewImageResponse(image))
}

// UploadImageFromURL handles image upload from a URL
func (c *ImageController) UploadImageFromURL(ctx *gin.Context) {
	var req struct {
		URL      string `json:"url" binding:"required"`
		FileName string `json:"file_name,omitempty"`
	}

	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	image, err := c.imageService.UploadImageFromURL(ctx.Request.Context(), req.URL, req.FileName)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewImageResponse(image))
}

func (c *ImageController) GetImageByID(ctx *gin.Context) {

	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	image, err := c.imageService.GetImageByID(ctx.Request.Context(), id)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
}

func (c *ImageController) GetAllImages(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetSource(common.CurrentService))
		return
	}
	req.Normalize()

	images, total, err := c.imageService.GetAllImages(ctx.Request.Context(), req.Page, req.Size)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	var imageResponses []dto.ImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, *dto.NewImageResponse(img))
	}

	page := dto.NewPaginationResponse(imageResponses, total, req)
	ctx.JSON(http.StatusOK, page)

}

func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	file, err := c.GetFile(ctx, "file")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	fileName := file.Filename
	f, ferr := file.Open()
	if ferr != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx).SetDetail(ferr.Error()).SetSource(common.CurrentService))
		return
	}
	defer f.Close()

	image, err := c.imageService.UpdateImageFromReader(ctx.Request.Context(), id, f, fileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
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

	image, err := c.imageService.UpdateImageFromURL(ctx.Request.Context(), id, req.URL, req.FileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
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
