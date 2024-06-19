package playlist

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
)

func (h *PlaylistHandler) GetUsersPlaylist(c *fiber.Ctx) error {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	pp, err := client.CurrentUsersPlaylists(c.Context(), spotify.Limit(pageSize), spotify.Offset(offset))
	if err != nil {
		return err
	}

	var ids []spotify.ID
	for _, p := range pp.Playlists {
		ids = append(ids, p.ID)
	}

	snapshotIDs, err := h.spotifyRepo.GetPlaylistSnapshotBulk(ids)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"playlists": pp,
		"saved":     snapshotIDs,
	})
}
