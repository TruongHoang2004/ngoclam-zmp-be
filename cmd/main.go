package main

import (
	"context"
	"log"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/module"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// fx.NopLogger, // Disable Fx's own logging
		module.ConfigModule(),
		module.DatabaseModule,
		module.ServerModule,
		module.RouterModule,
		module.ProductModule,
		module.ImageModule,
		module.FolderModule,

		// Hook cho log start/stop toàn app
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: OnStart,
				OnStop:  OnStop,
			})
		}),
	)

	app.Run()
}

func OnStart(ctx context.Context) error {
	log.Println("Application is starting...")
	return nil
}

func OnStop(ctx context.Context) error {
	log.Println("Application is stopping...")
	return nil
}
