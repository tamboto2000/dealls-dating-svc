package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamboto2000/dealls-dating-svc/internal/app"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
)

func main() {
	logger.Info("starting service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	app, err := app.NewApp(ctx, cfg)

	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := app.Start(); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("service started!")

	sigc := make(chan os.Signal, 1)
	signal.Notify(
		sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sigc

	logger.Info("stopping service...")
	cancel()
	if err := app.Stop(); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("service stopped")

	// TODO: recovery from panic
}
