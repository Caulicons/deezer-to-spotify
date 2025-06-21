package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

type GetSpotifyProfile struct {
	token *entities.SpotifyToken
}

func NewGetSpotifyProfile(token *entities.SpotifyToken) *GetSpotifyProfile {

	return &GetSpotifyProfile{
		token,
	}
}

func (u *GetSpotifyProfile) Execute() (profile map[string]any, erro *response.Err) {

	// Create request to Spotify API
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return profile, response.NewBadREquest(fmt.Sprintf("Error creating request: %v", err))
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token.AccessToken))

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return profile, response.NewBadREquest(fmt.Sprintf("Error fetching user profile: %v", err))
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return profile, response.NewBadREquest(fmt.Sprintf("Failed to get user profile: status code %d", resp.StatusCode))
	}

	// Parse and forward the response
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return profile, response.NewBadREquest(fmt.Sprintf("Error parsing user profile: %v", err))
	}

	return
}
