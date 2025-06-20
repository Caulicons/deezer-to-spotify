package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
)

type SpotifyAuthHandler struct {
	Auth *entities.SpotifyAuth
}

func NewSpotifyAuthHandler(auth *entities.SpotifyAuth) *SpotifyAuthHandler {

	return &SpotifyAuthHandler{
		Auth: auth,
	}
}

// authorizationURL generates the Spotify authorization URL
func (s *SpotifyAuthHandler) authorizationURL() string {
	u, _ := url.Parse("https://accounts.spotify.com/authorize")
	q := u.Query()

	q.Set("client_id", s.Auth.ClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", s.Auth.RedirectURI)
	q.Set("state", s.Auth.State)
	q.Set("scope", strings.Join(s.Auth.Scopes, " "))

	u.RawQuery = q.Encode()
	return u.String()
}

// RedirectToSpotifyAuth redirects the user to the Spotify authorization page
func (s *SpotifyAuthHandler) RedirectToSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.authorizationURL(), http.StatusFound)
}

// CallBack handles the Spotify OAuth callback and exchanges the code for an access token
func (s *SpotifyAuthHandler) CallBack(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	// Verify state to prevent CSRF attacks
	if state != s.Auth.State {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	// Prepare token request
	tokenURL := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.Auth.RedirectURI)

	// Create the request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.Auth.ClientID, s.Auth.ClientSecret)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Token request failed: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&s.Auth.Token); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	// Return the token data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&s.Auth.Token)
}
