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

	savedId, _ := h.spotifyRepo.GetPlaylistSnapshot(c.Context(), p.ID)
	if savedId != p.SnapshotID {
		//changed

		h.logger.Info("playlist changed", zap.String("id", p.ID.String()))

		h.spotifyRepo.CreateOrUpdatePlaylist(c.Context(), models.SpotifyPlaylist{
			ID:         p.ID.String(),
			SnapshotID: p.SnapshotID,
			Name:       p.Name,
			Images:     models.SpotifyImageFromList(p.Images),
		})

		go func() {
			tracks, err := getAllItemsForPlaylist(client, p.ID)
			if err != nil {
				h.logger.Error("failed to get playlist items", zap.Error(err))

				return
			}

			h.spotifyRepo.Refresh(context.Background(), p.ID.String(), tracks)
		}()
	}

	return nil
}

// get all tracks
func getAllItemsForPlaylist(client *spotify.Client, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	var items []spotify.PlaylistItem
	pageSize := 50
	// get the first 100 tracks
	playlist, err := client.GetPlaylistItems(context.Background(), playlistID, spotify.Limit(pageSize))
	if err != nil {
		return nil, err
	}
	items = append(items, playlist.Items...)
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

		pageNo++
		time.Sleep(100 * time.Millisecond)
	}

	return items, nil
}
