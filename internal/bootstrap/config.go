package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"go.uber.org/fx"
)

func ConfigModule() fx.Option {
	return fx.Provide(config.InitConfig)
}
