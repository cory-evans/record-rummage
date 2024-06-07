package apifx

import (
	"github.com/cory-evans/record-rummage/internal/api"
	"github.com/cory-evans/record-rummage/internal/api/handlers/auth"
	"github.com/cory-evans/record-rummage/internal/api/handlers/playlist"
	"github.com/cory-evans/record-rummage/internal/api/handlers/track"
	"github.com/cory-evans/record-rummage/internal/configfx"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/cory-evans/record-rummage/pkg/spotifyfx"
	"go.uber.org/fx"
)

var Module = fx.Module("api",
	configfx.Module,
	spotifyfx.Module,
	fx.Provide(
		api.NewApi,
		middleware.NewSpotifyClient,
		AsRoute(auth.NewAuthHandler),
		AsRoute(track.NewTrackHandler),
		AsRoute(playlist.NewPlaylistHandler),
	),
	fx.Invoke(func(api *api.Api) {}),
)

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(api.ApiRoute)),
		fx.ResultTags(`group:"api-routes"`),
	)
}
