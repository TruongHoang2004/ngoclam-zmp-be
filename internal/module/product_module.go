package module

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repository"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/controller"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/service"
	"github.com/gin-gonic/gin"

	"go.uber.org/fx"
)

var ProductModule = fx.Options(
	fx.Provide(repository.NewProductRepository),
	fx.Provide(service.NewProductService),
	fx.Provide(controller.NewProductController),

	fx.Invoke(func(h *controller.ProductController, g *gin.RouterGroup) {
		h.RegisterRoutes(g)
	}),
)
