package playlist

import (
	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/database"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PlaylistHandler struct {
	config *config.ApplicationConfig
	logger *zap.Logger
	router *fiber.App

	spotifyRepo   *database.SpotifyRepository
	spotifyAuth   *spotifyauth.Authenticator
	spotifyClient *middleware.SpotifyClient
}

type playlistHandlerParams struct {
	fx.In

	Config            *config.ApplicationConfig
	Logger            *zap.Logger
	SpotifyRepository *database.SpotifyRepository
	SpotifyAuth       *spotifyauth.Authenticator
	SpotifyClient     *middleware.SpotifyClient
}

func NewPlaylistHandler(p playlistHandlerParams) *PlaylistHandler {
	x := &PlaylistHandler{
		router:        fiber.New(),
		logger:        p.Logger,
		config:        p.Config,
		spotifyRepo:   p.SpotifyRepository,
		spotifyAuth:   p.SpotifyAuth,
		spotifyClient: p.SpotifyClient,
	}

	x.router.Post("/refresh", x.Refresh)

	return x
}

func (h *PlaylistHandler) Pattern() string {
	return "/playlist"
}

func (h *PlaylistHandler) Handler() *fiber.App {
	return h.router
}
