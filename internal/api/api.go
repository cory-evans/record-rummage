package api

import (
	"context"

	"github.com/cory-evans/record-rummage/internal/config"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/fx"
)

type Api struct {
	server *fiber.App
}

type apiParams struct {
	fx.In

	LC     fx.Lifecycle
	Routes []ApiRoute `group:"api-routes"`
	Config *config.ApplicationConfig
}

func NewApi(p apiParams) *Api {

	mux := fiber.New()

	mux.Use(cors.New())
	mux.Use(logger.New())
	mux.Use(middleware.NewSessionCookieMiddleware(p.Config))

	for _, route := range p.Routes {
		mux.Mount(route.Pattern(), route.Handler())
	}

	x := &Api{
		server: mux,
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go x.server.Listen(":8080")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return x.server.ShutdownWithContext(ctx)
		},
	})

	return x
}
