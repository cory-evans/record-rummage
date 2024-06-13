package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/zmb3/spotify/v2"
)

type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func (i *SpotifyImage) Scan(val interface{}) error {
	return scanJson(val, i)
}

func (i *SpotifyImage) Value() (driver.Value, error) {
	return json.Marshal(i)
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

type SpotifyImageArray []SpotifyImage

func (a *SpotifyImageArray) Scan(val interface{}) error {
	return scanJson(val, a)
}
