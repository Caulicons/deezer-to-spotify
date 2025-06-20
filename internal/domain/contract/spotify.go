package contract

import "net/http"

type ISpotifyAuth interface {
	RedirectToSpotifyAuth(w http.ResponseWriter, r *http.Request)
	CallBack(w http.ResponseWriter, r *http.Request)
}
type ISpotify interface {
	Me(w http.ResponseWriter, r *http.Request)
	CreatePlaylist(w http.ResponseWriter, r *http.Request)
	SearchAll(w http.ResponseWriter, r *http.Request)
	AddTracksToPlaylist(w http.ResponseWriter, r *http.Request)
}
