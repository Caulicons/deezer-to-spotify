package config

import (
	"net/http"

	dependencies "github.com/caulicons/deezer-to-spotify/cmd/api/depedencies"
)

func MapRoutes(mux *http.ServeMux, depend dependencies.Dependencies) {

	// Spotify
	mux.HandleFunc("/spotify/auth", depend.SpotifyHandler.Auth.RedirectToSpotifyAuth)
	mux.HandleFunc("/spotify/callback", depend.SpotifyHandler.Auth.CallBack)

	mux.HandleFunc("/spotify/me", depend.SpotifyHandler.Resource.Me)
	mux.HandleFunc("/spotify/playlist/{name}", depend.SpotifyHandler.Resource.CreatePlaylist)
	mux.HandleFunc("/spotify/tracks/search", depend.SpotifyHandler.Resource.SearchAll)
	mux.HandleFunc("/spotify/playlist/{id}/add", depend.SpotifyHandler.Resource.AddTracksToPlaylist)
	mux.HandleFunc("/spotify/me/playlists", depend.SpotifyHandler.Resource.GetAllUserPlaylists)
}
