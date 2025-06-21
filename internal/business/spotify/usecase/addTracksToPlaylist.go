package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	"github.com/caulicons/deezer-to-spotify/pkg/jsonUtils"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

// SpotifyTrackFound represents a track found in Spotify
type SpotifyTrackFound struct {
	ID    string `json:"id"`
	URI   string `json:"uri"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Isrc  string `json:"isrc"`
}

type SpotifyAddTracksToPlaylist struct {
	playlistID string
	token      *entities.SpotifyToken
}

func NewSpotifyAddTracksToPlaylist(playlistID string, token *entities.SpotifyToken) *SpotifyAddTracksToPlaylist {

	return &SpotifyAddTracksToPlaylist{
		playlistID,
		token,
	}
}

func (u *SpotifyAddTracksToPlaylist) Execute() (res map[string]any, erro *response.Err) {

	// Read track URIs from file
	tracks, err := jsonUtils.Read[SpotifyTrackFound]("spotify/track_uri.json")
	if err != nil {
		return res, response.NewInternalErr(fmt.Sprintf("Error reading track URIs: %v", err))
	}

	// Prepare track URIs for the request
	var trackURIs []string
	for _, track := range tracks {
		trackURIs = append(trackURIs, track.URI)
	}

	// Spotify API limits to 100 tracks per request, so we need to batch
	for i := 0; i < len(trackURIs); i += 100 {
		end := min(i+100, len(trackURIs))

		batchURIs := trackURIs[i:end]

		// Prepare API endpoint
		apiURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", u.playlistID)

		// Prepare request body
		requestBody := map[string]interface{}{
			"uris": batchURIs,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return res, response.NewInternalErr(fmt.Sprintf("Failed to create request body: %v", err))

		}

		// Create the request
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			return res, response.NewInternalErr(fmt.Sprintf("Failed to create request: %v", err))

		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token.AccessToken))

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return res, response.NewInternalErr(fmt.Sprintf("Failed to add tracks to playlist: %v", err))

		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			var errorResponse map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errorResponse)
			return res, response.NewInternalErr(fmt.Sprintf("Failed to add tracks. Status: %d, Error: %v", resp.StatusCode, errorResponse))
		}
	}

	res = map[string]any{
		"status":       "completed",
		"playlist_url": fmt.Sprintf("https://open.spotify.com/playlist/%s", u.playlistID),
	}
	return
}
