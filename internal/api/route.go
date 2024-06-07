package api

import (
	"github.com/gofiber/fiber/v2"
)

type ApiRoute interface {
	Pattern() string
	Handler() *fiber.App
}
