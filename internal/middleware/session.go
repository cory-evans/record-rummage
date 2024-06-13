package middleware

import (
	"time"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type SessionCookie struct {
	jwt.RegisteredClaims

	SpotifyUserID string        `json:"spotify_user_id"`
	SpotifyToken  *oauth2.Token `json:"spotify_token"`
}

const cookieKey = "session"

func NewSessionCookieMiddleware(appConfig *config.ApplicationConfig) fiber.Handler {

	signingKeyAsBytes := []byte(appConfig.JWTSigningKey)

	return func(c *fiber.Ctx) error {
		sessionCookieString := c.Cookies(cookieKey, "")

		if sessionCookieString == "" {
			return c.Next()
		}

		token, err := jwt.ParseWithClaims(
			sessionCookieString,
			&SessionCookie{},
			func(token *jwt.Token) (interface{}, error) {
				return signingKeyAsBytes, nil
			},
		)

		if err != nil {
			return c.Next()
		}

		claims, ok := token.Claims.(*SessionCookie)
		if !ok {
			return c.Next()
		}

		c.Locals("session", claims)

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
