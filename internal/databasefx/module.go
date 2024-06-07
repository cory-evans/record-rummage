package databasefx

import (
	"github.com/cory-evans/record-rummage/internal/database"
	"go.uber.org/fx"
)

var Module = fx.Module("databasefx",
	fx.Provide(
		database.NewDatabase,
		database.NewSpotifyRepository,
	),
)
