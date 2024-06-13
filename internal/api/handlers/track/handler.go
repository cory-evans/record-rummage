package track

import (
	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/database"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
)

type TrackHandler struct {
	config *config.ApplicationConfig

	router *fiber.App

	spotifyRepo   *database.SpotifyRepository
	spotifyAuth   *spotifyauth.Authenticator
	spotifyClient *middleware.SpotifyClient
}

type trackHandlerParams struct {
	fx.In

	Config *config.ApplicationConfig

	SpotifyRepository *database.SpotifyRepository
	SpotifyAuth       *spotifyauth.Authenticator
	SpotifyClient     *middleware.SpotifyClient
}

func NewTrackHandler(p trackHandlerParams) *TrackHandler {
	x := &TrackHandler{
		router:        fiber.New(),
		config:        p.Config,
		spotifyRepo:   p.SpotifyRepository,
		spotifyAuth:   p.SpotifyAuth,
		spotifyClient: p.SpotifyClient,
	}

	x.router.Use(middleware.NewSpotifyTokenMiddleware(middleware.SpotifyTokenMiddlewareConfig{}, p.Config))

	x.router.Get("/currently-playing", x.CurrentlyPlaying)
	x.router.Put("/playback", x.Playback)
	x.router.Post("/next", x.Next)
	x.router.Post("/previous", x.Previous)
	x.router.Get("/reveal", x.Reveal)

	return x
}

func (h *TrackHandler) Pattern() string {
	return "/track"
}

func (h *TrackHandler) Handler() *fiber.App {
	return h.router
}
