package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/services"
	"go.uber.org/fx"
)

func BuildService() fx.Option {
	return fx.Provide(
		services.NewBaseService,
		services.NewFolderService,
		services.NewImageService,
		services.NewProductService,
	)
}
