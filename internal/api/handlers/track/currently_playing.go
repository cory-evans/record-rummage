package track

import "github.com/gofiber/fiber/v2"

func (h *TrackHandler) CurrentlyPlaying(c *fiber.Ctx) error {

	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	playing, err := client.PlayerCurrentlyPlaying(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(playing)
}
