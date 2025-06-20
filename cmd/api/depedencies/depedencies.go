package dependencies

import (
	"github.com/caulicons/deezer-to-spotify/internal/business/spotify/usecase"
	"github.com/caulicons/deezer-to-spotify/internal/domain/contract"
	"github.com/caulicons/deezer-to-spotify/internal/infra/http/handler"
)

type Dependencies struct {
	SpotifyHandler struct {
		Auth     contract.ISpotifyAuth
		Resource contract.ISpotify
	}
}

func BuildDependencies() (depend Dependencies, err error) {

	// Spotify
	spotifyAuthConfig := usecase.NewSpotifyAuth("playlist-modify-public", "playlist-modify-private")
	spotifyAuth := handler.NewSpotifyAuthHandler(spotifyAuthConfig)
	depend.SpotifyHandler.Auth = spotifyAuth
	depend.SpotifyHandler.Resource = handler.NewSpotifyHandler(&spotifyAuth.Auth.Token)

	return
}
