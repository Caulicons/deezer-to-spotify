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

type SpotifyAddTracksToLoveSongs struct {
	token *entities.SpotifyToken
}

func NewSpotifyAddTrackToLoveSongs(token *entities.SpotifyToken) *SpotifyAddTracksToLoveSongs {

	return &SpotifyAddTracksToLoveSongs{
		token,
	}
}

func (u *SpotifyAddTracksToLoveSongs) Execute() (res map[string]any, erro *response.Err) {

	tracks, err := jsonUtils.Read[SpotifyTrackFound]("spotify/track_uri.json")
	if err != nil {
		return res, response.NewInternalErr(fmt.Sprintf("Error reading track URIs: %v", err))
	}

	// Prepare track URIs for the request
	var tracksID []string
	for _, track := range tracks {
		tracksID = append(tracksID, track.ID)
	}

	// Spotify Limit to add 50 Tracks per request at the Love Song, so we need to batch
	for i := 0; i < len(tracksID); i += 50 {
		end := min(i+50, len(tracksID))
		batchIDs := tracksID[i:end]

		// Prepare the request
		apiURL := "https://api.spotify.com/v1/me/tracks"
		requestBody := map[string]any{
			"ids": batchIDs,
		}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return res, response.NewInternalErr(fmt.Sprintf("Failed to create request body: %v", err))
		}

		// Create the request
		req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonBody))
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
			var errorResponse map[string]any
			json.NewDecoder(resp.Body).Decode(&errorResponse)
			return res, response.NewInternalErr(fmt.Sprintf("Failed to add tracks. Status: %d, Error: %v", resp.StatusCode, errorResponse))
		}
	}

	res = map[string]any{
		"Loved Song URL": "https://open.spotify.com/collection/tracks",
	}
	return
}
