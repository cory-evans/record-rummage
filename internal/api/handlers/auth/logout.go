package auth

import (
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Logout(c *fiber.Ctx) error {

	middleware.ClearSession(c)

	return c.Redirect("/")
}
