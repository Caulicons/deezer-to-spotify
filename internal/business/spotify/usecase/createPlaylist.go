package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

type SpotifyCreatePlaylist struct {
	name  string
	token *entities.SpotifyToken
}

func NewSpotifyCreatePlaylist(name string, token *entities.SpotifyToken) *SpotifyCreatePlaylist {

	return &SpotifyCreatePlaylist{
		name,
		token,
	}
}

func (u *SpotifyCreatePlaylist) Execute() (playlistID string, erro *response.Err) {

	// // Get user ID (needed to create a playlist)
	// userID, err := u.getUserID()
	// if err != nil {
	// 	return playlistID, response.NewInternalErr("Failed to get user ID: " + err.Error())
	// }

	// Create the playlist
	playlistURL := fmt.Sprintf("https://api.spotify.com/v1/me/playlists")

	// Prepare the request body
	requestBody := map[string]interface{}{
		"name":        u.name,
		"description": "Created via Music App",
		"public":      true,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return playlistID, response.NewInternalErr("Failed to create request body")
	}

	// Create the request
	req, err := http.NewRequest("POST", playlistURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return playlistID, response.NewInternalErr("Failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token.AccessToken))

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return playlistID, response.NewInternalErr("Failed to create playlist: " + err.Error())
	}
	defer resp.Body.Close()

	// FIX: Maybe i need add a better struct to the response later
	// Parse and return the response
	var playlistResponse map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&playlistResponse); err != nil {
		return playlistID, response.NewInternalErr("Failed to parse playlist response")
	}

	playlistID, ok := playlistResponse["id"].(string)
	if !ok {
		return playlistID, response.NewInternalErr("Failed to get playlist ID from response")
	}

	return
}
