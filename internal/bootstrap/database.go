package bootstrap

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/infrastructure/database"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var DatabaseModule = fx.Module("database",
	fx.Provide(database.NewDatabase),

	// đảm bảo close connection khi app stop
	fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				sqlDB, err := db.DB()
				if err != nil {
					return err
				}
				return sqlDB.Close()
			},
		})
	}),
)
