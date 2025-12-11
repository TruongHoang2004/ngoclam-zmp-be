package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/persistence/repositories"
	"go.uber.org/fx"
)

func BuildRepository() fx.Option {
	return fx.Provide(
		repositories.NewBaseRepository,
		repositories.NewProductRepository,
		repositories.NewImageRepository,
		repositories.NewFolderRepository,
		repositories.NewCategoryRepository,
		repositories.NewOrderRepository,
	)
}
