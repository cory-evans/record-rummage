package main

import (
	"github.com/cory-evans/record-rummage/internal/apifx"
	"github.com/cory-evans/record-rummage/internal/databasefx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		databasefx.Module,
		apifx.Module,
		fx.Provide(
			zap.NewProduction,
		)).Run()
}
