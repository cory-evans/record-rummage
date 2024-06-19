package middleware

import (
	"net/http"
	"time"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

type SpotifyClient struct {
	spotifyauth *spotifyauth.Authenticator
	spotifyRepo *database.SpotifyRepository
}

type spotifyClientParams struct {
	fx.In

	SpotifyAuth *spotifyauth.Authenticator
	SpotifyRepo *database.SpotifyRepository
}

func NewSpotifyClient(p spotifyClientParams) *SpotifyClient {
	return &SpotifyClient{
		spotifyauth: p.SpotifyAuth,
		spotifyRepo: p.SpotifyRepo,
	}
}

func (s *SpotifyClient) ForUser(c *fiber.Ctx) (*spotify.Client, error) {
	session := GetSession(c)

	if session == nil {
		return nil, fiber.ErrUnauthorized
	}

	httpClient := s.spotifyauth.Client(c.Context(), session.SpotifyToken)

	client := spotify.New(
		httpClient,
	)

	return client, nil
}

func (s *SpotifyClient) ForToken(c *fiber.Ctx, token *oauth2.Token) (*spotify.Client, error) {
	client := spotify.New(s.spotifyauth.Client(c.Context(), token))

	return client, nil
}

type SpotifyTokenMiddlewareConfig struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(*fiber.Ctx) bool
}

func NewSpotifyTokenMiddleware(config SpotifyTokenMiddlewareConfig, appConfig *config.ApplicationConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		session := GetSession(c)

		if (session == nil) || (session.SpotifyToken == nil) {
			return fiber.ErrUnauthorized
		}

		oauthToken := session.SpotifyToken

		// check if expired, request a new token
		if oauthToken.Expiry.UTC().Unix() < time.Now().UTC().Unix() {
			auth := spotifyauth.New(
				spotifyauth.WithRedirectURL(appConfig.SpotifyConfig.RedirectURI),
				spotifyauth.WithClientID(appConfig.SpotifyConfig.ClientID),
				spotifyauth.WithClientSecret(appConfig.SpotifyConfig.ClientSecret),
				spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadEmail),
			)
			client := spotify.New(auth.Client(c.Context(), oauthToken))
			newToken, err := client.Token()
			if err != nil {
				return err
			}

			SetSession(appConfig, c, newToken, session.SpotifyUserID)
		}

		err := c.Next()

		if err == nil {
			return nil
		}

		if sErr, ok := err.(spotify.Error); ok {
			if sErr.Status == http.StatusUnauthorized {
				// clear token
				ClearSession(c)

				return fiber.ErrUnauthorized
			}
		}

		return err

	}
}
