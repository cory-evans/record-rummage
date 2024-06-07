package auth

import (
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	cookie := middleware.GetSession(c)

	return c.JSON(cookie)
}
