package configfx

import (
	"github.com/cory-evans/record-rummage/internal/config"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		config.NewApplicationConfig,
	),
)
