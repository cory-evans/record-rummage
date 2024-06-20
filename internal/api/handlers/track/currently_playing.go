package track

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/cory-evans/record-rummage/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

func (h *TrackHandler) CurrentlyPlaying(c *fiber.Ctx) error {
	client, err := h.spotifyClient.ForUser(c)
	if err != nil {
		return err
	}

	playing, err := client.PlayerCurrentlyPlaying(c.Context())

	if err != nil {
		return err
	}

	go checkCurrentlyPlaying(h.logger, h.spotifyRepo, client, playing.PlaybackContext, playing.Item.ID)

	return c.JSON(playing)
}

func checkCurrentlyPlaying(logger *zap.Logger, repo *database.SpotifyRepository, client *spotify.Client, playbackContext spotify.PlaybackContext, trackId spotify.ID) {
	if playbackContext.Type != "playlist" {
		return
	}

	playlistURIParts := strings.Split(string(playbackContext.URI), ":")
	playlistId := playlistURIParts[len(playlistURIParts)-1]

	_, err := repo.GetPlaylistSnapshot(spotify.ID(playlistId))

	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("playlist not found", zap.String("playlist_id", playlistId))
		return
	}

	if err != nil {
		logger.Error("GetPlaylistSnapshot failed", zap.Error(err))
		return
	}

	addedBy, err := repo.GetAddedBy(playlistId, trackId.String())
	if err != nil {
		logger.Error("GetAddedBy failed", zap.Error(err))
		return
	}

	if len(addedBy) != 0 {
		logger.Debug("track already added", zap.String("track_id", trackId.String()))
		return
	}

	ctx := context.Background()

	// at this point we don't know who added the track, so we need to find out
	playlist, err := client.GetPlaylist(ctx, spotify.ID(playlistId))
	if err != nil {
		logger.Error("client.GetPlaylist failed", zap.Error(err))
		return
	}

	totalTracks := int(playlist.Tracks.Total)

	found := false

	limit := 50
	offset := totalTracks - limit

	// work backwards through the playlist to find the track
	for !found {
		page, err := client.GetPlaylistItems(
			ctx,
			spotify.ID(playlistId),
			spotify.Limit(limit),
			spotify.Offset(offset),
		)

		if err != nil {
			logger.Error("client.GetPlaylistItems failed", zap.Error(err))
			return
		}

		offset -= limit

		// merge into database
		err = repo.Refresh(playlistId, page.Items, false)
		if err != nil {
			logger.Error("repo.Refresh failed", zap.Error(err))
			return
		}

		for _, item := range page.Items {
			if item.Track.Track.ID == trackId {
				found = true
				break
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
