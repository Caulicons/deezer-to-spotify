package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// https://api.deezer.com/track/
func GetTrackInfo[T any](url string, id int) (data T, err error) {

	formattedURL := fmt.Sprintf("%s/%d", url, id)
	res, err := http.Get(formattedURL)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Fatalf("Error Decoding the request GET request: %v", err)
		return
	}

	return
}
