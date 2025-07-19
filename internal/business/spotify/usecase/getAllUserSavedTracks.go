package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (u *GetAllUserSavedTracks) Execute(token *entities.SpotifyToken) ([]entities.SpotifyPlaylist, *response.Err) {
	var allTracks []entities.SpotifyPlaylist
	nextURL := savedTrackDefaultURL

	var count = 1
	for nextURL != "" {
		fmt.Println(count, count)
		req, err := http.NewRequest(http.MethodGet, nextURL, nil)
		if err != nil {
			return nil, response.NewInternalErr(err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, response.NewInternalErr(err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, response.NewInternalErr(fmt.Sprintf("failed to get playlists, status code: %d", resp.StatusCode))
		}

		var FavoriteSongsResponse struct {
			Items []struct {
				Track entities.SpotifyPlaylist `json:"track"`
			} `json:"items"`
			Next string `json:"next"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&FavoriteSongsResponse); err != nil {
			return nil, response.NewInternalErr(fmt.Sprintf("Error decodend the playlistReponse : %v", err))
		}

		// Extract tracks from items and append to allTracks
		for _, item := range FavoriteSongsResponse.Items {
			allTracks = append(allTracks, item.Track)
		}
		nextURL = FavoriteSongsResponse.Next

		fmt.Println(allTracks)

		fmt.Println(count)
		count++
	}

	jsonUtils.Write(allTracks, "spotify/favorite_tracks.json")
	return nil, nil
}
