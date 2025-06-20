package config

import (
	"net/http"

	dependencies "github.com/caulicons/deezer-to-spotify/cmd/api/depedencies"
	"github.com/caulicons/deezer-to-spotify/internal/infra/http/middleware"
)

func MapRoutes(mux *http.ServeMux, depend dependencies.Dependencies) {

	// Spotify
	mux.HandleFunc("/spotify/auth", depend.SpotifyHandler.Auth.RedirectToSpotifyAuth)
	mux.HandleFunc("/spotify/callback", depend.SpotifyHandler.Auth.CallBack)

	// mux.Handle("/playlist/{name}", middleware.AuthMiddleware())
	mux.Handle("/spotify/me", middleware.SpotifyAuth(http.HandlerFunc(depend.SpotifyHandler.Resource.Me)))
	mux.HandleFunc("/spotify/playlist/{name}", depend.SpotifyHandler.Resource.CreatePlaylist)
	mux.HandleFunc("/spotify/tracks/search", depend.SpotifyHandler.Resource.SearchAll)
	mux.HandleFunc("/spotify/playlist/{id}/add", depend.SpotifyHandler.Resource.AddTracksToPlaylist)
}
