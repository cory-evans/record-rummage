package database

import (
	"context"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type newDatabaseParams struct {
	fx.In
	Log       *zap.Logger
	LC        fx.Lifecycle
	AppConfig *config.ApplicationConfig
}

func NewDatabase(params newDatabaseParams) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), params.AppConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}

	params.LC.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close(context.Background())
		},
	})

	return conn, nil
}
