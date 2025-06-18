package main

import (
	"fmt"

	"github.com/caulicons/deezer-to-spotify/internal/deezer/usecases"
)

func main() {

	err := run()

	if err != nil {
		fmt.Println(err)
	}
}

var baseDeezerURL = "https://api.deezer.com/user/1237568626/tracks?index=0"

func run() error {

	url := baseDeezerURL
	for {
		_, next, err := usecases.GetPlaylistTracks(url)
		url = next
		if err != nil {
			return err
		}

		if url == "" {
			break
		}

	}
	return nil
}
