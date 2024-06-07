package auth

import (
	"github.com/cory-evans/record-rummage/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	err := h.spotifyRepo.CheckLoginState(
		c.Context(),
		state,
		true,
	)
	if err != nil {
		return err
	}

	token, err := h.spotifyAuth.Exchange(c.Context(), code)
	if err != nil {
		return err
	}

	client, err := h.spotifyClient.ForToken(c, token)
	if err != nil {
		return err
	}

	user, err := client.CurrentUser(c.Context())
	if err != nil {
		return err
	}

	middleware.SetSession(h.config, c, token, user.ID)

	return c.Redirect("http://localhost:4200/")
}
