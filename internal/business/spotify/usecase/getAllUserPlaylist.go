package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caulicons/deezer-to-spotify/internal/constants"
	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

type GetAllUserPlaylistsSpotify struct {
	httpClient *http.Client
}

func NewGetAllUserPlaylistsSpotify(httpClient *http.Client) *GetAllUserPlaylistsSpotify {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &GetAllUserPlaylistsSpotify{
		httpClient: httpClient,
	}
}

func (u *GetAllUserPlaylistsSpotify) Execute(token *entities.SpotifyToken) ([]entities.SpotifyPlaylist, *response.Err) {
	var allPlaylists []entities.SpotifyPlaylist
	nextURL := constants.SpotifyPlaylistEndpoint
	limit := 50

	for nextURL != "" {
		url := fmt.Sprintf("%s?limit=%d", nextURL, limit)
		if nextURL != constants.SpotifyPlaylistEndpoint {
			url = nextURL // If it's a pagination URL, use it directly
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, response.NewInternalErr(err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		resp, err := u.httpClient.Do(req)
		if err != nil {
			return nil, response.NewInternalErr(err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, response.NewInternalErr(fmt.Sprintf("failed to get playlists, status code: %d", resp.StatusCode))
		}

		var playlistsResponse struct {
			Items []entities.SpotifyPlaylist `json:"items"`
			Next  string                     `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&playlistsResponse); err != nil {
			return nil, response.NewInternalErr(fmt.Sprintf("Error decodend the playlistReponse : %v", err))

		}

		allPlaylists = append(allPlaylists, playlistsResponse.Items...)
		nextURL = playlistsResponse.Next

		if nextURL == "" {
			break // No more pages to fetch
		}
	}

	return allPlaylists, nil
}
