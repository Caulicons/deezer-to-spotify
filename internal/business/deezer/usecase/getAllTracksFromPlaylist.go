package usecase

import "log"

func GetAllTracksFromPlaylist[T any](url string) (tracks []T, err error) {

	for {
		data, _, next, err := GetTracksFromPlaylist[T](url)
		if err != nil {
			log.Fatalf("Error Getting All Trackings: %v", err)
			return tracks, err
		}

		if next == "" {
			return tracks, nil
		}

		tracks = append(tracks, data...)
		url = next
	}
}
