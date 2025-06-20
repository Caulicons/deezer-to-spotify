package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	json1 "github.com/caulicons/deezer-to-spotify/pkg/utils/json"
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
	// Check if we have an access token
	if err := s.checkToken(w); err != nil {
		return
	}

	// Create request to Spotify API
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user profile: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to get user profile: status code %d", resp.StatusCode), resp.StatusCode)
		return
	}

	// Parse and forward the response
	var userProfile map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing user profile: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the user profile data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userProfile)
}

// Search All Tracks
func (s *SpotifyHandler) SearchAll(w http.ResponseWriter, r *http.Request) {
	var count = 1

	if err := s.checkToken(w); err != nil {
		return
	}

	tracks, err := json1.Read[entities.DeezerTrackInfo]("deezer/track_info.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting the tracks: %v", err), http.StatusInternalServerError)
		return
	}

	type SpotifyTrackFound struct {
		ID    string `json:"id"`
		URI   string `json:"uri"`
		Name  string `json:"name"`
		Title string `json:"title"`
		Isrc  string `json:"isrc"`
	}

	type TrackNotFound struct {
		Title string `json:"title"`
		Isrc  string `json:"isrc"`
	}

	tracksFound := []SpotifyTrackFound{}
	tracksNotFound := []TrackNotFound{}

	for _, track := range tracks {
		// if count == 15 {
		// 	http.Error(w, "finishing tests", 800)
		// 	break
		// }

		// Construct search query using track name and ISRC
		searchQuery := url.QueryEscape(fmt.Sprintf("%s isrc:%s", track.Title, track.Isrc))
		searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", searchQuery)

		// Create request
		req, err := http.NewRequest("GET", searchURL, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating search request: %v", err), http.StatusInternalServerError)
			return
		}

		// Set authorization header
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error searching track: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Parse response
		var searchResult struct {
			Tracks struct {
				Items []struct {
					ID   string `json:"id"`
					URI  string `json:"uri"`
					Name string `json:"name"`
				} `json:"items"`
			} `json:"tracks"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing search results: %v", err), http.StatusInternalServerError)
			return
		}

		// Check if we found a match
		if len(searchResult.Tracks.Items) > 0 {
			fmt.Printf("%d - Found track: %s (ID: %s, URI: %s)\n",
				count,
				searchResult.Tracks.Items[0].Name,
				searchResult.Tracks.Items[0].ID,
				searchResult.Tracks.Items[0].URI)

			// Add to found tracks
			tracksFound = append(tracksFound, SpotifyTrackFound{
				ID:    searchResult.Tracks.Items[0].ID,
				URI:   searchResult.Tracks.Items[0].URI,
				Name:  searchResult.Tracks.Items[0].Name,
				Title: track.Title,
				Isrc:  track.Isrc,
			})
		} else {
			fmt.Printf("%d - No match found for track: %s by %s\n", count, track.Title, track.Isrc)

			//Add to not found tracks
			tracksNotFound = append(tracksNotFound, TrackNotFound{
				Title: track.Title,
				Isrc:  track.Isrc,
			})
		}
		count++
	}

	// Write results to files
	if err := json1.Write(tracksFound, "spotify/track_uri.json"); err != nil {
		http.Error(w, fmt.Sprintf("Error writing found tracks: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json1.Write(tracksNotFound, "spotify/tracks_not_found.json"); err != nil {
		http.Error(w, fmt.Sprintf("Error writing not found tracks: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}

// CreatePlaylist creates a new Spotify playlist with the name provided in URL path
func (s *SpotifyHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {

	if err := s.checkToken(w); err != nil {
		return
	}

	playlistName := r.PathValue("name")
	fmt.Println(playlistName)

	// Get user ID (needed to create a playlist)
	userID, err := s.getUserID()
	if err != nil {
		http.Error(w, "Failed to get user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the playlist
	playlistURL := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID)

	// Prepare the request body
	requestBody := map[string]interface{}{
		"name":        playlistName,
		"description": "Created via Music App",
		"public":      false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Failed to create request body", http.StatusInternalServerError)
		return
	}

	// Create the request
	req, err := http.NewRequest("POST", playlistURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to create playlist: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse and return the response
	var playlistResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&playlistResponse); err != nil {
		http.Error(w, "Failed to parse playlist response", http.StatusInternalServerError)
		return
	}

	// Return the playlist data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlistResponse)
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

	// Read track URIs from file
	tracks, err := json1.Read[SpotifyTrackFound]("spotify/track_uri.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading track URIs: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare track URIs for the request
	var trackURIs []string
	for _, track := range tracks {
		trackURIs = append(trackURIs, track.URI)
	}

	// Spotify API limits to 100 tracks per request, so we need to batch
	for i := 0; i < len(trackURIs); i += 100 {
		end := i + 100
		if end > len(trackURIs) {
			end = len(trackURIs)
		}

		batchURIs := trackURIs[i:end]

		// Prepare API endpoint
		apiURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

		// Prepare request body
		requestBody := map[string]interface{}{
			"uris": batchURIs,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request body: %v", err), http.StatusInternalServerError)
			return
		}

		// Create the request
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add tracks to playlist: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			var errorResponse map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errorResponse)
			http.Error(w, fmt.Sprintf("Failed to add tracks. Status: %d, Error: %v", resp.StatusCode, errorResponse), http.StatusInternalServerError)
			return
		}
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Added %d tracks to playlist", len(trackURIs)),
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

// getUserID retrieves the current user's Spotify ID
func (s *SpotifyHandler) getUserID() (string, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user profile: %d", resp.StatusCode)
	}

	var userProfile struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		return "", err
	}

	return userProfile.ID, nil
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
