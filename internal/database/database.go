package database

import (
	"context"

	"github.com/cory-evans/record-rummage/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type newDatabaseParams struct {
	fx.In
	Log       *zap.Logger
	LC        fx.Lifecycle
	AppConfig *config.ApplicationConfig
}

func NewDatabase(params newDatabaseParams) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("pgx", params.AppConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}

	params.LC.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})

	return conn, nil
}
