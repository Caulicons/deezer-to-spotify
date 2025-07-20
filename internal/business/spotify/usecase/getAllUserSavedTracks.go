package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/constants"
	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	"github.com/caulicons/deezer-to-spotify/pkg/jsonUtils"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

const savedTrackDefaultURL = "https://api.spotify.com/v1/me/tracks?limit=50"

type GetAllUserSavedTracks struct {
}

func NewGetAllUserSavedTracks() *GetAllUserSavedTracks {

	return &GetAllUserSavedTracks{}
}

func (u *GetAllUserSavedTracks) Execute(token *entities.SpotifyToken) *response.Err {
	var allTracks []entities.SpotifyPlaylist
	nextURL := savedTrackDefaultURL

	for nextURL != "" {
		req, err := http.NewRequest(http.MethodGet, nextURL, nil)
		if err != nil {
			return response.NewInternalErr(err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return response.NewInternalErr(err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return response.NewInternalErr(fmt.Sprintf("failed to get playlists, status code: %d", resp.StatusCode))
		}

		var FavoriteSongsResponse struct {
			Items []struct {
				Track entities.SpotifyPlaylist `json:"track"`
			} `json:"items"`
			Next string `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&FavoriteSongsResponse); err != nil {
			return response.NewInternalErr(fmt.Sprintf("Error decodend the playlistReponse : %v", err))
		}

		// Extract tracks from items and append to allTracks
		for _, item := range FavoriteSongsResponse.Items {
			allTracks = append(allTracks, item.Track)
		}
		nextURL = FavoriteSongsResponse.Next
	}

	jsonUtils.Write(allTracks, constants.SpotifyUserFavoriteTracksFile)

	return nil
}
