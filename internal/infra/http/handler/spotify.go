package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/business/spotify/usecase"
	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
)

type SpotifyHandler struct {
	token *entities.SpotifyToken
}

func NewSpotifyHandler(token *entities.SpotifyToken) *SpotifyHandler {
	return &SpotifyHandler{
		token,
	}
}

func (s *SpotifyHandler) Me(w http.ResponseWriter, r *http.Request) {

	if err := s.checkToken(w); err != nil {
		return
	}

	usecase := usecase.NewGetSpotifyProfile(s.token)
	userProfile, err := usecase.Execute()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	// Return the user profile data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userProfile)
}

// Search All Tracks
func (s *SpotifyHandler) SearchAll(w http.ResponseWriter, r *http.Request) {

	usecase := usecase.NewSpotifySearchAllTracks(s.token)
	res, err := usecase.Execute()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// CreatePlaylist creates a new Spotify playlist with the name provided in URL path
func (s *SpotifyHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {

	if err := s.checkToken(w); err != nil {
		return
	}

	playlistName := r.PathValue("name")
	fmt.Println(playlistName)

	usecase := usecase.NewSpotifyCreatePlaylist(playlistName, s.token)
	playlistID, err := usecase.Execute()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	// Return the playlist data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     fmt.Sprintf("Playlist with name %s was created with ID: %s", playlistName, playlistID),
		"status":      "completed",
		"playlist_id": playlistID,
	})
}

// Add the tracks to a specific playlist
func (s *SpotifyHandler) AddTracksToPlaylist(w http.ResponseWriter, r *http.Request) {

	if err := s.checkToken(w); err != nil {
		return
	}

	// Get playlist ID from path parameter
	playlistID := r.PathValue("id")
	if playlistID == "" {
		http.Error(w, "Playlist ID is required", http.StatusBadRequest)
		return
	}

	usecase := usecase.NewSpotifyAddTracksToPlaylist(playlistID, s.token)
	res, err := usecase.Execute()

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Added %d tracks to playlist", len(res)),
	})
}

// SpotifyTrackFound represents a track found in Spotify
type SpotifyTrackFound struct {
	ID    string `json:"id"`
	URI   string `json:"uri"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Isrc  string `json:"isrc"`
}

// getAuthConfig extracts and validates authentication information from the request
func (s *SpotifyHandler) checkToken(w http.ResponseWriter) error {
	token := s.token.AccessToken
	if token == "" {
		message := "missing authorization token go to and authorizate http://127.0.0.1:8080/spotify/auth"
		http.Error(w, message, http.StatusInternalServerError)
		return errors.New(message)
	}

	return nil
}
