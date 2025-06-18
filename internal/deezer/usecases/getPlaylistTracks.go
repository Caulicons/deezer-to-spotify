package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PlaylistTrackRes struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type PlaylistRes struct {
	Data  []PlaylistTrackRes `json:"data"`
	Total int                `json:"total"`
	Prev  string             `json:"prev"`
	Next  string             `json:"next"`
}

var count = 1

func GetPlaylistTracks(url string) (prev string, next string, err error) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
		return prev, next, err
	}
	defer res.Body.Close()

	var playlistRes PlaylistRes
	err = json.NewDecoder(res.Body).Decode(&playlistRes)
	if err != nil {
		log.Fatalf("Error Converting the request GET request: %v", err)
		return prev, next, err
	}

	fmt.Println(playlistRes)
	fmt.Println("Count : ", count)
	count++

	prev = playlistRes.Prev
	next = playlistRes.Next
	return prev, next, err
}
