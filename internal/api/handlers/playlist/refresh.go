package playlist

import (
	"context"
	"time"

	"github.com/cory-evans/record-rummage/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

func (h *PlaylistHandler) Refresh(c *fiber.Ctx) error {

	playlistId := c.Query("id")

	if playlistId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing playlist id")
	}

	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	p, err := client.GetPlaylist(c.Context(), spotify.ID(playlistId))
	if err != nil {
		return err
	}

	savedId, _ := h.spotifyRepo.GetPlaylistSnapshot(p.ID)
	if savedId != p.SnapshotID {
		//changed

		h.logger.Info("playlist changed", zap.String("id", p.ID.String()))

		h.spotifyRepo.CreateOrUpdatePlaylist(models.SpotifyPlaylist{
			ID:         p.ID.String(),
			SnapshotID: p.SnapshotID,
			Name:       p.Name,
			Images:     models.SpotifyImageFromList(p.Images),
		})

		h.CurrentlyWorkingPlaylists[playlistId] = 0

		go func() {
			tracks, err := h.getAllItemsForPlaylist(client, p.ID)
			if err != nil {
				h.logger.Error("failed to get playlist items", zap.Error(err))

				return
			}

			h.spotifyRepo.Refresh(p.ID.String(), tracks, true)
		}()
	}

	return nil
}

// get all tracks
func (h *PlaylistHandler) getAllItemsForPlaylist(client *spotify.Client, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	var items []spotify.PlaylistItem
	pageSize := 50
	// get the first 100 tracks
	playlist, err := client.GetPlaylistItems(context.Background(), playlistID, spotify.Limit(pageSize))
	if err != nil {
		return nil, err
	}

	playlistIDString := string(playlistID)

	items = append(items, playlist.Items...)
	h.CurrentlyWorkingPlaylists[playlistIDString] = float64(len(items)) / float64(playlist.Total)

	// get the rest of the tracks
	pageNo := 1
	for playlist.Next != "" {
		playlist, err = client.GetPlaylistItems(
			context.Background(),
			playlistID,
			spotify.Offset(pageNo*pageSize),
			spotify.Limit(pageSize),
		)

		if err != nil {
			return nil, err
		}

		items = append(items, playlist.Items...)

		h.CurrentlyWorkingPlaylists[playlistIDString] = float64(len(items)) / float64(playlist.Total)

		pageNo++
		time.Sleep(100 * time.Millisecond)
	}

	delete(h.CurrentlyWorkingPlaylists, playlistIDString)

	return items, nil
}

func (h *PlaylistHandler) RefreshProgress(c *fiber.Ctx) error {

	playlistId := c.Query("id")

	if playlistId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing playlist id")
	}

	var progress float64 = 1

	if p, ok := h.CurrentlyWorkingPlaylists[playlistId]; ok {
		progress = p
	}

	return c.JSON(fiber.Map{
		"progress":   progress,
		"playlistID": playlistId,
	})
}
