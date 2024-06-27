package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Api struct {
	server *fiber.App
}

type apiParams struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	Routes []ApiRoute `group:"api-routes"`
	Config *config.ApplicationConfig

	SpotifyAuth *spotifyauth.Authenticator
}

func NewApi(p apiParams) *Api {
	mux := fiber.New()

	apiGroup := mux.Group("/api")

	apiGroup.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	apiGroup.Use(
		middleware.NewSessionCookieMiddleware(
			p.Config,
			p.Logger,
			p.SpotifyAuth,
			func(c *fiber.Ctx) bool {
				return strings.HasPrefix(c.Path(), "/api/auth/") && c.Path() != "/api/auth/me"
			},
		),
	)

	for _, route := range p.Routes {
		apiGroup.Mount(route.Pattern(), route.Handler())
	}

	if p.Config.IsDev {
		mux.Get("/*", func(c *fiber.Ctx) error {
			err := proxy.Do(c, "http://localhost:4200"+c.Path())
			if err != nil {
				p.Logger.Error("proxy error", zap.Error(err))
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			return nil
		})

	} else {
		// make sure path exists
		_, err := os.Stat(path.Join("wwwroot", "index.html"))
		if errors.Is(err, os.ErrNotExist) {
			p.Logger.Error("index.html does not exist", zap.Error(err))
			os.Exit(1)
		}

		mux.Use("/", filesystem.New(filesystem.Config{
			Root: NewSpaFileSystem(
				http.Dir("wwwroot"),
			),
			Index: "index.html",
		}))
	}

	x := &Api{
		server: mux,
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := ":80"
			if p.Config.IsDev {
				port = ":8080"
			}

			go x.server.Listen(port)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return x.server.ShutdownWithContext(ctx)
		},
	})

	return x
}
