package middleware

import (
	"time"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type SessionCookie struct {
	jwt.RegisteredClaims

	SpotifyUserID string        `json:"spotify_user_id"`
	SpotifyToken  *oauth2.Token `json:"spotify_token"`
}

const cookieKey = "session"

func NewSessionCookieMiddleware(
	appConfig *config.ApplicationConfig,
	logger *zap.Logger,
	spotifyauth *spotifyauth.Authenticator,
	skip func(*fiber.Ctx) bool,
) fiber.Handler {

	signingKeyAsBytes := []byte(appConfig.JWTSigningKey)

	return func(c *fiber.Ctx) error {

		if skip != nil && skip(c) {
			return c.Next()
		}

		sessionCookieString := c.Cookies(cookieKey, "")

		if sessionCookieString == "" {
			return fiber.ErrUnauthorized
		}

		token, err := jwt.ParseWithClaims(
			sessionCookieString,
			&SessionCookie{},
			func(token *jwt.Token) (interface{}, error) {
				return signingKeyAsBytes, nil
			},
		)

		if err != nil {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(*SessionCookie)
		if !ok {
			return fiber.ErrUnauthorized
		}

		if claims.SpotifyToken.Expiry.Before(time.Now().UTC()) {

			newToken, err := spotifyauth.RefreshToken(c.Context(), claims.SpotifyToken)
			if err != nil {
				logger.Error("failed to refresh token", zap.Error(err))
				return fiber.ErrUnauthorized
			}

			err = SetSession(appConfig, c, newToken, claims.SpotifyUserID)
			if err != nil {
				logger.Error("failed to set session", zap.Error(err))
				return fiber.ErrUnauthorized
			}
		} else {
			c.Locals("session", claims)
		}

		return c.Next()
	}
}

func SetSession(appConfig *config.ApplicationConfig, c *fiber.Ctx, spotifyToken *oauth2.Token, spotifyUserID string) error {
	claims := &SessionCookie{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "record-rummage",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       util.GenerateRandomString(32),
			Subject:  spotifyUserID,
		},
		SpotifyUserID: spotifyUserID,
		SpotifyToken:  spotifyToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(appConfig.JWTSigningKey))
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieKey,
		Value:    ss,
		HTTPOnly: true,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 30),
	})

	c.Locals("session", claims)

	return nil
}

func ClearSession(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieKey,
		Value:    "",
		HTTPOnly: true,
	})
}

func GetSession(c *fiber.Ctx) *SessionCookie {
	session := c.Locals("session")
	if session == nil {
		return nil
	}

	return session.(*SessionCookie)
}
