package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/caulicons/deezer-to-spotify/cmd/shell/config"
	deezerUS "github.com/caulicons/deezer-to-spotify/internal/business/deezer/usecase"
	spotifyUS "github.com/caulicons/deezer-to-spotify/internal/business/spotify/usecase"
	"github.com/caulicons/deezer-to-spotify/internal/constants"
	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	"github.com/caulicons/deezer-to-spotify/pkg/jsonUtils"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
func run() error {
	// Display welcome message
	Welcome()

	// Step 1: Ask user for Deezer playlist
	err := DeezerMenu()
	if err != nil {
		return fmt.Errorf("failed to process Deezer Menu: %v", err)
	}

	// Step 2: Spotify Authentication
	spotifyToken, err := SpotifyAuthenticator()
	if err != nil {
		return err
	}

	// Step 3: Main Menu Loop
	err = SpotifyMenu(spotifyToken)
	if err != nil {
		return err
	}

	return nil
}

func Welcome() {

	fmt.Println("\n======================================================")
	fmt.Println("🎵 Welcome to Deezer to Spotify Playlist Converter 🎵")
	fmt.Println("This tool helps you transfer your favorite tracks from Deezer to Spotify")
	fmt.Printf("======================================================\n\n")
}

func DeezerOptions() error {
	// URL formats examples
	// Favorite Playlist = "https://api.deezer.com/user/{your_playlist_ID}/"
	// Any other Public Playlist = "https://api.deezer.com/playlist/{your_playlist_ID}/tracks"
	// Documentation : https://developers.deezer.com/api/playlist

	var formattedURL string
outer:
	for {
		fmt.Println("\n🎵 Select your Deezer playlist type:")
		fmt.Println("1 - Your 'Loved' Playlist")
		fmt.Println("2 - Any other Public Playlist")
		fmt.Println("3 - Canceled")
		fmt.Print("> ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {

		case 1:
			fmt.Printf("Enter your Deezer user ID: ")
			fmt.Println("❔➜ Go to your Deezer profile in the URL you will see some random numbers, that's it.")

			var userID string
			fmt.Print("> ")
			fmt.Scanln(&userID)

			formattedURL = fmt.Sprintf("https://api.deezer.com/user/%s/tracks", userID)
			break outer

		case 2:
			fmt.Printf("Enter the Deezer playlist ID: ")
			fmt.Println("❔➜ When you access your Deezer playlist through Browser, a random number always appears in the URL, the ID is that.")

			var playlistID string
			fmt.Print("> ")
			fmt.Scanln(&playlistID)
			formattedURL = fmt.Sprintf("https://api.deezer.com/playlist/%s/tracks", playlistID)
			break outer

		case 3:
			return fmt.Errorf("user canceled operation")
		}
	}

	// Get Tracks from the Deezer Playlist
	tracks, err := deezerUS.GetAllTracksFromPlaylistDeezer[entities.DeezerPlaylistTrackData](formattedURL)
	if err != nil {
		return err
	}

	// Get Tracks Info from the Deezer Playlist, this is to get the ISRC
	trackInfo, err := deezerUS.GetTrackInfoBatchGetID[entities.DeezerPlaylistTrackData, entities.DeezerTrackInfo]("https://api.deezer.com/track", tracks,
		func(dptd entities.DeezerPlaylistTrackData) int {
			return dptd.ID
		},
		func(dptd entities.DeezerPlaylistTrackData) string {
			return dptd.Title
		},
	)
	if err != nil {
		return err
	}

	jsonUtils.Write(trackInfo, constants.DeezerTracksFile)

	return nil
}

func DeezerMenu() error {

	fmt.Println("First we will need your playlist from the Deezer")
	err := DeezerOptions()
	if err != nil {
		return fmt.Errorf("failed to process Deezer playlist: %v", err)
	}

	return nil
}

func SpotifyAuthenticator() (*entities.SpotifyToken, error) {
	// Step 2: Spotify Authentication
	fmt.Println("\n🔐 You need to authenticate with Spotify. Press ENTER to continue...")
	fmt.Scanln()

	spotifyToken, err := config.StartSpotifyAuthServer()
	if err != nil {

		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if spotifyToken == nil || spotifyToken.AccessToken == "" {
		return nil, fmt.Errorf("auth complete but token is empty")
	}

	fmt.Println("✅ Spotify Auth Successful!")

	return spotifyToken, nil
}

func SpotifyMenu(token *entities.SpotifyToken) error {
	httpClient := &http.Client{}
	for {
		fmt.Println("\n🎛️  What do you want to do?")
		fmt.Println("1. List all Spotify playlists")
		fmt.Println("2. Create a new Spotify playlist")
		fmt.Println("3. Reset All Love song tracks")
		fmt.Println("4. Exit")
		fmt.Print("> ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// List all playlists
			allPlaylists, err := spotifyUS.NewGetAllUserPlaylistsSpotify(httpClient).Execute(token)
			allPlaylists = append([]entities.SpotifyPlaylist{{Name: "♥️ Love Songs"}}, allPlaylists...)
			if err != nil {
				fmt.Println("❌ Error fetching playlists:", err.Message)
				continue
			}

			fmt.Println("\nAvailable playlists:")
			for i, playlist := range allPlaylists {
				fmt.Printf("%d. %s \n", i+1, playlist.Name)
			}

			fmt.Println("\nChoose a playlist by number to add Deezer tracks to it:")
			var idx int
			fmt.Print("> ")
			fmt.Scanln(&idx)

			if idx <= 0 || idx > len(allPlaylists) {
				fmt.Println("❌ Invalid selection.")
				continue
			}

			selected := allPlaylists[idx-1]

			// Get tracks Info (search and map) in Spotify
			res, err := spotifyUS.NewSpotifySearchAllTracks(token).Execute()
			if err != nil {
				fmt.Println("❌ Error searching tracks:", err.Message)
				continue
			}

			fmt.Println("\n⏳ Tracks Result: ", selected.Name)
			for k, v := range res {
				fmt.Println(k, ":", v)
			}

			// Confirm with the user before adding tracks
			fmt.Println("\nDo you want to add these tracks to the selected Spotify playlist?")
			fmt.Println("1. Yes")
			fmt.Println("2. No")
			fmt.Print("> ")

			var confirmation int
			fmt.Scanln(&confirmation)

			if confirmation != 1 {
				fmt.Println("Operation canceled. No tracks were added to the playlist.")
				continue
			}

			fmt.Println("Adding tracks to playlist... 💨")

			if selected.Name == "♥️ Love Songs" {
				res, err = spotifyUS.NewSpotifyAddTrackToLoveSongs(token).Execute()
				if err != nil {
					fmt.Println("❌ Error adding tracks:", err.Message)
					continue
				}
			} else {
				res, err = spotifyUS.NewSpotifyAddTracksToPlaylist(selected.ID, token).Execute()
				if err != nil {
					fmt.Println("❌ Error adding tracks:", err.Message)
					continue
				}
			}

			fmt.Println("\n✅ Tracks added to", selected.Name)
			for k, v := range res {
				fmt.Println(k, ":", v)
			}
		case 2:
			// Create new playlist
			fmt.Println("Enter name for the new playlist:")
			fmt.Print("> ")

			// Use bufio.NewReader to handle spaces in input
			reader := bufio.NewReader(os.Stdin)
			playlistName, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Error reading input:", err)
				continue
			}

			// Trim newline characters
			playlistName = strings.TrimSpace(playlistName)

			if playlistName == "" {
				fmt.Println("❌ Playlist name cannot be empty")
				continue
			}

			_, erro := spotifyUS.NewSpotifyCreatePlaylist(playlistName, token).Execute()
			if erro != nil {
				fmt.Println("❌ Error creating playlist:", erro.Message)
				continue
			}

			fmt.Println("✅ Playlist created! ")

		case 3:
			fmt.Println("choice 4")
			_, err := spotifyUS.NewGetAllUserSavedTracks().Execute(token)
			if err != nil {
				fmt.Println("❌ Error Get tracks:", err.Message)
			}

			_, err = spotifyUS.NewDeleteAllUserSavedTracks().Execute(token)
			if err != nil {
				fmt.Println("❌ Error Get tracks:", err.Message)
			}

		case 4:
			fmt.Println("👋 Goodbye!")
			return nil

		default:
			fmt.Println("❌ Invalid option.")
		}
	}
}
