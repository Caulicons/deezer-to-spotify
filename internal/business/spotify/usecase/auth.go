package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
)

// NewSpotifyAuth creates a new SpotifyAuth instance
func NewSpotifyAuth() *entities.SpotifyAuth {
	return &entities.SpotifyAuth{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URL"),
		State:        generateRandomString(16),
		Scopes:       []string{"playlist-modify-public", "playlist-modify-private", "playlist-read-private", "user-library-modify", "user-library-read"},
	}
}

// generateRandomString creates a random string for the state parameter
func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:length]
}
