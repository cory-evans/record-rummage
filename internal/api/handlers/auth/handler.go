package auth

import (
	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/database"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
)

type AuthHandler struct {
	config *config.ApplicationConfig

	router *fiber.App

	spotifyRepo   *database.SpotifyRepository
	spotifyAuth   *spotifyauth.Authenticator
	spotifyClient *middleware.SpotifyClient
}

type authHandlerParams struct {
	fx.In

	Config *config.ApplicationConfig

	SpotifyRepository *database.SpotifyRepository
	SpotifyAuth       *spotifyauth.Authenticator
	SpotifyClient     *middleware.SpotifyClient
}

func NewAuthHandler(p authHandlerParams) *AuthHandler {
	x := &AuthHandler{
		router:        fiber.New(),
		config:        p.Config,
		spotifyRepo:   p.SpotifyRepository,
		spotifyAuth:   p.SpotifyAuth,
		spotifyClient: p.SpotifyClient,
	}

	x.router.Get("/login", x.Login)
	x.router.Get("/logout", x.Logout)
	x.router.Get("/callback", x.Callback)
	x.router.Get("/me", x.GetMe)

	return x
}
func (h *AuthHandler) Pattern() string {
	return "/auth"
}

func (h *AuthHandler) Handler() *fiber.App {
	return h.router
}
