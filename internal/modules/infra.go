package modules

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/infra"
	"github.com/tamboto2000/dealls-dating-svc/pkg/cache"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

type Infra struct {
	db *sqli.DB
	ch *cache.Cache
}

func NewInfra(ctx context.Context, cfg config.Config) (*Infra, error) {
	// database
	db, err := infra.InitDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// cache
	ch, err := infra.InitCache(ctx, cfg)
	if err != nil {
		return nil, err
	}

	in := Infra{
		db: db,
		ch: ch,
	}

	return &in, nil
}

func (in *Infra) DB() *sqli.DB {
	return in.db
}

func (in *Infra) Cache() *cache.Cache {
	return in.ch
}
