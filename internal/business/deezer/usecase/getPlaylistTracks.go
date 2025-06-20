package usecase

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
)

func GetTracksFromPlaylist[T any](url string) (data []T, prev string, next string, err error) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
		return
	}
	defer res.Body.Close()

	var playlist entities.PaginationRes[T]
	err = json.NewDecoder(res.Body).Decode(&playlist)
	if err != nil {
		log.Fatalf("Error Converting the request GET request: %v", err)
		return
	}
	return playlist.Data, playlist.Prev, playlist.Next, err
}
