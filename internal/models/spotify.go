package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/zmb3/spotify/v2"
)

type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func (i *SpotifyImage) Scan(value interface{}) error {
	src, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type")
	}

	return json.Unmarshal(src, i)
}

func SpotifyImageFromList(src []spotify.Image) []SpotifyImage {
	var images = make([]SpotifyImage, 0, len(src))
	for _, i := range src {
		images = append(images, SpotifyImage{
			URL:    i.URL,
			Height: int(i.Height),
			Width:  int(i.Width),
		})
	}
	return images
}

type SpotifyUser struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"display_name"`
	Images      []SpotifyImage `json:"images"`
}

type SpotifyArtist struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []SpotifyImage `json:"images"`
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
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []SpotifyImage `json:"images"`
}

type SpotifyTrack struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
	Album   SpotifyAlbum    `json:"album"`
}

type SpotifyPlaylist struct {
	ID         string         `json:"id"`
	SnapshotID string         `json:"snapshot_id"`
	Name       string         `json:"name"`
	Images     []SpotifyImage `json:"images"`
}

type SpotifyPlaylistTrack struct {
	PlaylistID string    `json:"playlist_id"`
	TrackID    string    `json:"track_id"`
	AddedAt    time.Time `json:"added_at"`
	AddedBy    string    `json:"added_by"`
}
