package track

import (
	"github.com/cory-evans/record-rummage/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
)

func (h *TrackHandler) Reveal(c *fiber.Ctx) error {
	trackId := c.Query("trackId")
	playlistId := c.Query("playlistId")

	addedBy, err := h.spotifyRepo.GetAddedBy(playlistId, trackId)
	if err != nil {
		return c.JSON(nil)
	}

	var users = make([]*models.SpotifyUser, 0)

	for _, id := range addedBy {
		user, err := h.spotifyRepo.GetUser(id)
		if err != nil {
			client, err := h.spotifyClient.ForUser(c)
			if err != nil {
				return err
			}

			spotifyUser, err := client.GetUsersPublicProfile(c.Context(), spotify.ID(id))
			if err != nil {
				return err
			}

			user = &models.SpotifyUser{
				ID:          spotifyUser.ID,
				DisplayName: spotifyUser.DisplayName,
				Images:      models.SpotifyImageFromList(spotifyUser.Images),
			}

			h.spotifyRepo.CreateOrUpdateUser(user)
		}

		users = append(users, user)
	}

	return c.JSON(users)
}
