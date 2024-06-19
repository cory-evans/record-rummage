package track

import "github.com/gofiber/fiber/v2"

func (h *TrackHandler) Playback(c *fiber.Ctx) error {

	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	body := struct {
		IsPlaying bool `json:"is_playing"`
	}{}

	err = c.BodyParser(&body)

	if err != nil {
		return err
	}

	if body.IsPlaying {
		err = client.Pause(c.Context())
	} else {
		err = client.Play(c.Context())
	}

	if err != nil {
		return err

	}

	return c.JSON(fiber.Map{})
}

func (h *TrackHandler) Next(c *fiber.Ctx) error {
	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	q, err := client.GetQueue(c.Context())
	if err != nil {
		return err
	}

	err = client.Next(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(q)
}
func (h *TrackHandler) Previous(c *fiber.Ctx) error {
	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	err = client.Previous(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{})
}
