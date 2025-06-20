package middleware

import (
	"context"
	"errors"
	"net/http"
)

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	Token string
}

// getAuthConfig extracts and validates authentication information from the request
func getAuthConfig(r *http.Request) (*AuthConfig, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("missing authorization token")
	}

	// Further validation can be added here

	return &AuthConfig{
		Token: token,
	}, nil
}

func SpotifyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate auth config/token
		authConfig, err := getAuthConfig(r)
		if err != nil {

			// Redirect to authentication endpoint
			http.Redirect(w, r, "http://127.0.0.1:8080/spotify/auth", http.StatusTemporaryRedirect)
			return
		}

		// Store in context and continue
		ctx := context.WithValue(r.Context(), "authConfig", authConfig)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
