package auth

import (
	"github.com/cory-evans/record-rummage/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	state := util.GenerateRandomString(16)

	err := h.spotifyRepo.CreateLoginState(
		state,
	)

	if err != nil {
		return err
	}

	url := h.spotifyAuth.AuthURL(state)
	return c.Redirect(url)
}
