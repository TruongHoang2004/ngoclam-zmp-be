package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/bootstrap"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/shopspring/decimal"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	defaultGracefulTimeout = 15 * time.Second
)

func init() {
	fmt.Println("Initializing application...")
	config.InitConfig()
	fmt.Println("Config is loaded")
	log.NewLogger()
	fmt.Println("Logger is initialized")
	decimal.MarshalJSONWithoutQuotes = true
}

func main() {

	logger := log.GetLogger().GetZap()
	logger.Debugf("App is running")

	app := fx.New(
		// fx.NopLogger, // Disable Fx's own logging
		fx.Provide(log.GetLogger().GetZap),
		bootstrap.BuildExtServicesModules(),
		bootstrap.BuildDatabase(),
		bootstrap.BuildRepository(),
		bootstrap.BuildService(),
		bootstrap.BuildController(),
		bootstrap.BuildValidator(),
		bootstrap.ServerModule,
		bootstrap.RouterModule,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		logger.Fatalf(err.Error())
	}

	interruptHandle(app, logger)
}

func OnStart(ctx context.Context) error {
	log.Info(ctx, "Application is starting...")
	return nil
}

func OnStop(ctx context.Context) error {
	log.Info(ctx, "Application is stopping...")
	return nil
}

func interruptHandle(app *fx.App, logger *zap.SugaredLogger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Debugf("Listening Signal...")
	s := <-c
	logger.Infof("Received signal: %s. Shutting down Server ...", s)

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Fatalf(err.Error())
	}
}
