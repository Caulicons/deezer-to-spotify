package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/caulicons/deezer-to-spotify/cmd/api/config"
	dependencies "github.com/caulicons/deezer-to-spotify/cmd/shell/depedencies"
	"github.com/caulicons/deezer-to-spotify/internal/domain/entities"
	"github.com/caulicons/deezer-to-spotify/internal/infra/http/handler"
)

type Server struct {
	Mux     *http.ServeMux
	address string
}

func NewServerConfig() *Server {

	mux := http.NewServeMux()
	port, has := os.LookupEnv("PORT")
	if !has {
		port = "8080"
	}

	return &Server{
		Mux:     mux,
		address: ":" + port,
	}
}

func StartSpotifyAuthServer() (chan bool, *entities.SpotifyToken) {
	// Prompt the user to open the Spotify authentication page in their browser
	fmt.Println("\n🌐 Opening Spotify authentication page in your browser...")

	// Initialize token to avoid nil reference
	token := &entities.SpotifyToken{}

	// Channel to receive authentication result
	tokenReadyChan := make(chan bool)
	var tokenMutex sync.Mutex

	// Start the auth server in a goroutine
	go func() {
		app, err := config.NewApplication()
		if err != nil {
			fmt.Println("Error creating application:", err)
			tokenReadyChan <- false
			return
		}

		depend, err := dependencies.BuildDependencies()
		if err != nil {
			fmt.Println("Error building dependencies:", err)
			tokenReadyChan <- false
			return
		}

		// Override the callback handler to capture token and signal completion
		authHandler := depend.SpotifyHandler.Auth.(*handler.SpotifyAuthHandler)
		authHandler.CallBackFunc = func(w http.ResponseWriter, r *http.Request) {
			// Copy the token after successful authentication
			tokenMutex.Lock()
			*token = authHandler.Auth.Token
			tokenMutex.Unlock()

			// Display success message on webpage
			w.Write([]byte(`{"status": "success", "message": "Authentication successful! You can close this window and return to the terminal."}`))

			// Signal that auth is complete
			tokenReadyChan <- true
		}

		MapRoutes(app.Server.Mux, depend)

		// Run the server with a shutdown mechanism
		server := &http.Server{
			Addr:    ":8080",
			Handler: app.Server.Mux,
		}

		go func() {
			<-tokenReadyChan
			// Wait a moment for the response to be sent
			time.Sleep(2 * time.Second)
			server.Shutdown(context.Background())
		}()

		showLog()

		server.ListenAndServe()
	}()

	return tokenReadyChan, token
}

func showLog() {
	// Try to automatically open the browser
	url := "http://127.0.0.1:8080/spotify/auth"
	var err error
	// Attempt to open browser based on operating system
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Could not open browser automatically:", err)
	}
	fmt.Println("📱 If it doesn't open automatically, please navigate to: http://127.0.0.1:8080/spotify/auth")

	// Try to open the browser automatically
	// URL is defined in the printed message above
	fmt.Printf("\n⏳ Waiting for authorization...\n\n\n")

}
