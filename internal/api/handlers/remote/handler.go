package remote

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RemoteHandler struct {
	router *fiber.App
	logger *zap.Logger
	pubsub *gochannel.GoChannel
}

type remoteHandlerParams struct {
	fx.In

	Logger *zap.Logger
}

func NewRemoteHandler(p remoteHandlerParams) *RemoteHandler {
	x := &RemoteHandler{
		router: fiber.New(),
		pubsub: gochannel.NewGoChannel(
			gochannel.Config{},
			watermill.NewCaptureLogger(),
		),
		logger: p.Logger,
	}

	x.router.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	x.router.Get("/ws/sub", websocket.New(func(c *websocket.Conn) {

		clientId := c.Query("clientID")
		if clientId == "" {
			x.logger.Sugar().Warn("no client id")
			return
		}

		x.logger.Sugar().Info("websocket connection")

		session := c.Locals("session")
		if session == nil {
			x.logger.Sugar().Warn("no session")
			return
		}

		s, ok := session.(*middleware.SessionCookie)
		if !ok {
			x.logger.Sugar().Warn("invalid session")
			return
		}

		sub, err := x.pubsub.Subscribe(context.Background(), "remote/"+s.SpotifyUserID)
		if err != nil {
			x.logger.Sugar().Error(err)
			return
		}

		for msg := range sub {
			if msg.Metadata.Get("client_id") == clientId {
				continue
			}

			err = c.WriteJSON(msg)
			if err != nil {
				x.logger.Sugar().Error(err)
				return
			}

			msg.Ack()
		}

	}))

	x.router.Post("/pub", func(c *fiber.Ctx) error {
		s := middleware.GetSession(c)
		msg := message.NewMessage(
			watermill.NewUUID(),
			[]byte(c.Query("message")),
		)

		msg.Metadata.Set("client_id", c.Query("clientID"))

		err := x.pubsub.Publish("remote/"+s.SpotifyUserID, msg)
		if err != nil {

			return err
		}

		return c.JSON(fiber.Map{})
	})

	return x
}

func (h *RemoteHandler) Pattern() string {
	return "/remote"
}

func (h *RemoteHandler) Handler() *fiber.App {
	return h.router
}
