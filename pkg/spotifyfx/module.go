package spotifyfx

import (
	"github.com/cory-evans/record-rummage/internal/config"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
)

var Module = fx.Module("spotify",
	fx.Provide(
		NewSpotifyAuth,
	),
)

func NewSpotifyAuth(appConfig *config.ApplicationConfig) *spotifyauth.Authenticator {
	return spotifyauth.New(
		spotifyauth.WithRedirectURL(appConfig.SpotifyConfig.RedirectURI),
		spotifyauth.WithClientID(appConfig.SpotifyConfig.ClientID),
		spotifyauth.WithClientSecret(appConfig.SpotifyConfig.ClientSecret),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
	)
}
