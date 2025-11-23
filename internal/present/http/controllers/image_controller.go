package controllers

import (
	"fmt"
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController(imageService *services.ImageService) *ImageController {
	return &ImageController{
		imageService: imageService,
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
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(common.BadRequest("File is required", err))
		return
	}

	// The original filename uploaded by the client
	fileName := file.Filename

	// Open the file to read its content
	f, err := file.Open()
	if err != nil {
		ctx.Error(common.BadRequest("Failed to open file", err))
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	image, err := c.imageService.UploadImageFromURL(ctx.Request.Context(), req.URL, req.FileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewImageResponse(image))
}

func (c *ImageController) GetImageByID(ctx *gin.Context) {

	idParam := ctx.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	image, err := c.imageService.GetImageByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
}

func (c *ImageController) GetAllImages(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid pagination query", err))
		return
	}
	req.Normalize()

	images, total, err := c.imageService.GetAllImages(ctx.Request.Context(), req.Page, req.Size)
	if err != nil {
		ctx.Error(err)
		return
	}

	var imageResponses []dto.ImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, *dto.NewImageResponse(img))
	}

	pag := dto.NewPaginationResponse(imageResponses, total, req)
	ctx.JSON(http.StatusOK, pag)

}

func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id := ctx.Param("id")
	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(common.BadRequest("File is required", err))
		return
	}

	fileName := file.Filename
	f, err := file.Open()
	if err != nil {
		ctx.Error(common.BadRequest("Failed to open file", err))
		return
	}
	defer f.Close()

	image, err := c.imageService.UpdateImageFromReader(ctx.Request.Context(), idUint, f, fileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
}

// UpdateImageFromURL handles image update from a URL
func (c *ImageController) UpdateImageFromURL(ctx *gin.Context) {
	id := ctx.Param("id")
	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	var req struct {
		URL      string `json:"url" binding:"required"`
		FileName string `json:"file_name,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(common.BadRequest("Invalid request body", err))
		return
	}

	image, err := c.imageService.UpdateImageFromURL(ctx.Request.Context(), idUint, req.URL, req.FileName)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewImageResponse(image))
}

func (c *ImageController) DeleteImage(ctx *gin.Context) {
	id := ctx.Param("id")
	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}

	err = c.imageService.DeleteImage(ctx.Request.Context(), idUint)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
