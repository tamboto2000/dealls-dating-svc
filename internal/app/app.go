package app

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/modules"
	"github.com/tamboto2000/dealls-dating-svc/internal/rest"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mux"
)

type App struct {
	in   *modules.Infra
	muxr *mux.Router
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	var app App

	// init snowid
	// TODO: node id to be retrieved from config
	if err := modules.InitSnowID(1); err != nil {
		return nil, err
	}

	// init inra
	in, err := modules.NewInfra(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// init components
	comps, err := modules.NewComponents(cfg)
	if err != nil {
		return nil, err
	}

	// init services
	svcs := modules.NewServices(ctx, cfg, in, comps)

	// init http router
	router := mux.NewRouter()
	rest.RegisterREST(cfg, router, svcs)

	app = App{
		in:   in,
		muxr: router,
	}

	return &app, nil
}

func (a *App) Start() error {
	go a.muxr.Run(":8000")

	return nil
}

func (a *App) Stop() error {
	// close database
	if err := a.in.DB().Close(); err != nil {
		logger.Error(err.Error())
	}

	// close cache
	if err := a.in.Cache().Close(); err != nil {
		logger.Error(err.Error())
	}

	return nil
}
