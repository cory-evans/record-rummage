package models

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

type SpotifyUser struct {
	ID          string            `json:"id" db:"id"`
	DisplayName string            `json:"display_name" db:"display_name"`
	Images      SpotifyImageArray `json:"images" db:"images"`
}

type SpotifyArtist struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Images SpotifyImageArray `json:"images"`
}

func SpotifyArtistFromSlice(src []spotify.SimpleArtist) []SpotifyArtist {
	var artists = make([]SpotifyArtist, 0, len(src))
	for _, i := range src {
		artists = append(artists, SpotifyArtist{
			ID:   i.ID.String(),
			Name: i.Name,
		})
	}
	return artists
}

type SpotifyAlbum struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Images SpotifyImageArray `json:"images"`
}

type SpotifyTrack struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
	Album   SpotifyAlbum    `json:"album"`
}

type SpotifyPlaylist struct {
	ID         string            `json:"id"`
	SnapshotID string            `json:"snapshot_id"`
	Name       string            `json:"name"`
	Images     SpotifyImageArray `json:"images"`
}

type SpotifyPlaylistTrack struct {
	PlaylistID string    `json:"playlist_id"`
	TrackID    string    `json:"track_id"`
	AddedAt    time.Time `json:"added_at"`
	AddedBy    string    `json:"added_by"`
}
