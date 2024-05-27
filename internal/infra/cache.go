package infra

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/pkg/cache"
)

func InitCache(ctx context.Context, cfg config.Config) (*cache.Cache, error) {
	c, err := cache.NewCache(cache.Options{
		Addr: cfg.Cache.Addr(),
	})

	if err != nil {
		return nil, err
	}

	return c, c.Ping(ctx)
}
