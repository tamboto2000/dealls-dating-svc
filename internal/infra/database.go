package infra

import (
	"context"
	"fmt"

	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

func InitDatabase(ctx context.Context, cfg config.Config) (*sqli.DB, error) {
	dbCfg := cfg.Database
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
	conn, err := sqli.NewPostgreConn(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	db := sqli.NewDB(conn)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
