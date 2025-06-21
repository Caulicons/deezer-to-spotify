package usecase

import "log"

func GetAllTracksFromPlaylistDeezer[T any](url string) (tracks []T, err error) {

	for {
		data, _, next, err := GetTracksFromPlaylist[T](url)
		if err != nil {
			log.Fatalf("Error Getting All Trackings: %v", err)
			return tracks, err
		}
		tracks = append(tracks, data...)

		if next == "" {
			return tracks, nil
		}

		url = next
	}
}
