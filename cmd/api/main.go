package main

import (
	"github.com/cory-evans/record-rummage/internal/apifx"
	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/databasefx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		databasefx.Module,
		apifx.Module,
		fx.Provide(
			NewLogger,
		)).Run()
}

func NewLogger(cfg *config.ApplicationConfig) *zap.Logger {
	var logger *zap.Logger

	if cfg.IsDev {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	return logger
}
