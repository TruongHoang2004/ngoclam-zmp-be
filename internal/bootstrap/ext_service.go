package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/sdk/imagekit"
	"go.uber.org/fx"
)

func BuildExtServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(imagekit.NewImageKitClient),
	)
}
