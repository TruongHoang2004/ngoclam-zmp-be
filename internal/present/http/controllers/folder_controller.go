package controllers

import (
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	httpCommon "github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/dto"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type FolderController struct {
	*baseController
	folderService *services.FolderService
}

func NewFolderController(
	base *baseController,
	folderService *services.FolderService,
) *FolderController {
	return &FolderController{
		baseController: base,
		folderService:  folderService,
	}
}

func (c *FolderController) RegisterRoutes(r *gin.RouterGroup) {
	folders := r.Group("/folders")
	{
		folders.POST("", c.CreateFolder)
		folders.GET("", c.ListFolders)
		folders.GET(":id", c.GetFolderByID)
		folders.PUT(":id", c.UpdateFolder)
		folders.DELETE(":id", c.DeleteFolder)
	}
}

func (c *FolderController) CreateFolder(ctx *gin.Context) {
	var req dto.CreateFolderRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	if err := c.folderService.CreateFolder(ctx, &req); err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusCreated, httpCommon.NewSuccessResponse(nil))
}

func (c *FolderController) GetFolderByID(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	folder, serviceErr := c.folderService.GetFolderByID(ctx, uint(id))
	if serviceErr != nil {
		ctx.JSON(serviceErr.HTTPStatus, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewFolderResponse(folder)))
}

func (c *FolderController) ListFolders(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx))
		return
	}
	req.Normalize()

	folders, total, err := c.folderService.ListFolders(ctx, req.Page, req.Size)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	var responses []dto.FolderResponse
	for _, f := range folders {
		responses = append(responses, *dto.NewFolderResponse(f))
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewPaginationResponse(responses, total, dto.PaginationRequest{Page: req.Page, Size: req.Size})))
}

func (c *FolderController) UpdateFolder(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	var req dto.UpdateFolderRequest
	if err := c.BindAndValidateRequest(ctx, &req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	updatedFolder, err := c.folderService.UpdateFolder(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(err.HTTPStatus, err)
		return
	}

	ctx.JSON(http.StatusOK, httpCommon.NewSuccessResponse(dto.NewFolderResponse(updatedFolder)))
}

func (c *FolderController) DeleteFolder(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	if err := c.folderService.DeleteFolder(ctx.Request.Context(), id); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Folder deleted"})
}
