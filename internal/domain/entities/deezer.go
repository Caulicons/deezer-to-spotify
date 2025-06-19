package entities

type DeezerPlaylistTrackData struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type DeezerTrackInfo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Isrc  string `json:"isrc"`
}
