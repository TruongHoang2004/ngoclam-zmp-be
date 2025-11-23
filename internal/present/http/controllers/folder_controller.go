package controllers

import (
	"net/http"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
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
	if err := c.BindAndValidateRequest(ctx, req); err != nil {
		c.ErrorData(ctx, err)
		return
	}

	f, err := c.folderService.CreateFolder(ctx.Request.Context(), req.Name, req.Description)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewFolderResponse(f))
}

func (c *FolderController) GetFolderByID(ctx *gin.Context) {
	id, err := c.GetUintParam(ctx, "id")
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	f, err := c.folderService.GetFolderByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewFolderResponse(f))
}

func (c *FolderController) ListFolders(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.ErrorData(ctx, common.ErrBadRequest(ctx))
		return
	}
	req.Normalize()

	folders, total, err := c.folderService.ListFolders(ctx.Request.Context(), req.Page, req.Size)
	if err != nil {
		c.ErrorData(ctx, err)
		return
	}

	var res []dto.FolderResponse
	for _, f := range folders {
		res = append(res, *dto.NewFolderResponse(f))
	}

	pag := dto.NewPaginationResponse(res, total, req)
	ctx.JSON(http.StatusOK, pag)
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

	updated, err := c.folderService.UpdateFolder(ctx.Request.Context(), id, &req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewFolderResponse(updated))
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
