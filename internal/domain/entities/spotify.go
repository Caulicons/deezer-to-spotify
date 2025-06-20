package entities

// SpotifyAuth handles Spotify authorization flow
type SpotifyAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	State        string
	Scopes       []string
	Token        SpotifyToken
}

type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
