package json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// For agreement (by my self only) all the json file are storage in the folder
func Write[T any](data T, path string) (err error) {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Create data directory if it doesn't exist
	dataDir := fmt.Sprintf("%s/data/", cwd)
	// if err := os.MkdirAll(dataDir, 0755); err != nil {
	// 	log.Fatalf("Error creating data directory: %v", err)
	// 	return err
	// }

	destPath := fmt.Sprintf("%s/%s", dataDir, path)
	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error opening the file %s: %v", destPath, err)
		return err
	}
	defer file.Close()

	jsonData, nil := json.Marshal(data)

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", destPath, err)
		return err
	}

	return nil
}

// type SpotifyTrackFound struct {
// 	ID    string `json:"id"`
// 	URI   string `json:"uri"`
// 	Name  string `json:"name"`
// 	Title string `json:"title"`
// 	Isrc  string `json:"isrc"`
// }

// type TrackNotFound struct {
// 	Title string `json:"title"`
// 	Isrc  string `json:"isrc"`
// }

// tracksFound := []SpotifyTrackFound{}
// tracksNotFound := []TrackNotFound{}

// Add to not found tracks
// 	tracksNotFound = append(tracksNotFound, TrackNotFound{
// 		Title: track.Title,
// 		Isrc:  track.Isrc,
// 	})
// }

// count++

// // Write results to files
// if err := jsonb.Write(tracksFound, "spotify/track_uri.json"); err != nil {
// 	http.Error(w, fmt.Sprintf("Error writing found tracks: %v", err), http.StatusInternalServerError)
// 	return
// }

// if err := jsonb.Write(tracksNotFound, "spotify/tracks_not_found.json"); err != nil {
// 	http.Error(w, fmt.Sprintf("Error writing not found tracks: %v", err), http.StatusInternalServerError)
// 	return
// }
