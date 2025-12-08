package bootstrap

import (
	"net/http"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/client/zalo"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/sdk/imagekit"
	"go.uber.org/fx"
)

func BuildExtServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(func() *http.Client {
			return &http.Client{
				Timeout: 10 * time.Second,
			}
		}),
		fx.Provide(imagekit.NewImageKitClient),
		fx.Provide(zalo.NewClient),
	)
}
