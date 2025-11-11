package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService *service.ImageService
}

func NewImageController(imageService *service.ImageService) *ImageController {
	return &ImageController{
		imageService: imageService,
	}
}

func (c *ImageController) RegisterRoutes(r *gin.RouterGroup) {
	images := r.Group("/images")
	{
		images.POST("/", c.UploadImage)
		images.GET("/:id", c.GetImageByID)
		images.GET("/", c.GetAllImages)
		images.PUT("/:id", c.UpdateImage)
		images.DELETE("/:id", c.DeleteImage)
	}
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithError(400, gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	// The original filename uploaded by the client
	fileName := file.Filename
	fmt.Println("Uploaded file name:", fileName)

	// Open the file to read its content
	f, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	image, _ := c.imageService.UploadImage(ctx.Request.Context(), fileName, fileData)
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
	page := ctx.DefaultQuery("page", "0")
	limit := ctx.DefaultQuery("limit", "10")

	var pageInt, limitInt int
	_, err := fmt.Sscanf(page, "%d", &pageInt)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid page format", err))
		return
	}
	_, err = fmt.Sscanf(limit, "%d", &limitInt)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid limit format", err))
		return
	}

	images, err := c.imageService.GetAllImages(ctx.Request.Context(), pageInt, limitInt)
	if err != nil {
		ctx.Error(err)
		return
	}

	var imageResponses []dto.ImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, *dto.NewImageResponse(img))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"images": imageResponses,
	})
	// return

}

func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id := ctx.Param("id")
	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		ctx.Error(common.BadRequest("Invalid ID format", err))
		return
	}
	image, err := c.imageService.UpdateImage(ctx.Request.Context(), idUint, "", nil)
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
	c.imageService.DeleteImage(ctx.Request.Context(), 0)
}
