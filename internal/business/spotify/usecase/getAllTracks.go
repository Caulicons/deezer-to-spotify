package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	"github.com/caulicons/deezer-to-spotify/pkg/jsonUtils"
	response "github.com/caulicons/deezer-to-spotify/pkg/reponse"
)

type SpotifySearchAllTracks struct {
	token *entities.SpotifyToken
}

func NewSpotifySearchAllTracks(token *entities.SpotifyToken) *SpotifySearchAllTracks {

	return &SpotifySearchAllTracks{
		token,
	}
}

func (u *SpotifySearchAllTracks) Execute() (res map[string]any, erro *response.Err) {

	var count = 1

	tracks, err := jsonUtils.Read[entities.DeezerTrackInfo]("deezer/track_info.json")
	if err != nil {
		return res, response.NewInternalErr(fmt.Sprintf("Error getting the tracks: %v", err))
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
			return res, response.NewInternalErr(fmt.Sprintf("Error creating search request: %v", err))
		}

		// Set authorization header
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token.AccessToken))

		// Execute request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return res, response.NewInternalErr(fmt.Sprintf("Error searching track: %v", err))
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
			return res, response.NewInternalErr(fmt.Sprintf("Error parsing search results: %v", err))
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
	if err := jsonUtils.Write(tracksFound, "spotify/track_uri.json"); err != nil {
		return res, response.NewInternalErr(fmt.Sprintf("Error parsing search results: %v", err))

	}

	if err := jsonUtils.Write(tracksNotFound, "spotify/tracks_not_found.json"); err != nil {
		return res, response.NewInternalErr(fmt.Sprintf("Error writing not found tracks: %v", err))

	}

	res = map[string]any{"status": "completed"}
	return
}
