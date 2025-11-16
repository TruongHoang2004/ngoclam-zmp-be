package module

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/controller"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/service"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/sdk/imagekit"
	"github.com/gin-gonic/gin"

	"go.uber.org/fx"
)

var ImageModule = fx.Options(
	fx.Provide(imagekit.NewImageKitClient),
	fx.Provide(repository.NewImageRepository),
	fx.Provide(service.NewImageService),
	fx.Provide(controller.NewImageController),

	fx.Invoke(func(h *controller.ImageController, g *gin.RouterGroup) {
		h.RegisterRoutes(g)
	}),
)
