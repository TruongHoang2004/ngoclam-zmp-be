package bootstrap

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/controllers"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/present/http/validator"

	"go.uber.org/fx"
)

func BuildController() fx.Option {
	return fx.Options(
		fx.Provide(controllers.NewBaseController),
		fx.Provide(controllers.NewFolderController),
		fx.Provide(controllers.NewProductController),
		fx.Provide(controllers.NewImageController),
		fx.Provide(controllers.NewCategoryController),
	)
}

func BuildValidator() fx.Option {
	return fx.Options(
		fx.Provide(validator.NewValidator),
		fx.Invoke(validator.RegisterDecimalTypeFunc),
		fx.Invoke(validator.RegisterValidations),
	)
}
