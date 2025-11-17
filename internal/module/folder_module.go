package module

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/controller"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var FolderModule = fx.Options(
	fx.Provide(repository.NewFolderRepository),
	fx.Provide(service.NewFolderService),
	fx.Provide(controller.NewFolderController),
	fx.Invoke(func(h *controller.FolderController, g *gin.RouterGroup) {
		h.RegisterRoutes(g)
	}),
)
