package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

func StartSpotifyAuthServer() (*entities.SpotifyToken, error) {
	// Prompt the user to open the Spotify authentication page in their browser
	fmt.Println("\n🌐 Opening Spotify authentication page in your browser...")

	// Initialize token to avoid nil reference
	token := &entities.SpotifyToken{}

	// Channel to receive authentication result
	tokenReady := make(chan bool)

	server, err := configSpotifyAuthServer(tokenReady, token)
	if err != nil {
		return nil, err
	}

	// Create context with timeout for server operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	go func() {
		showLog()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error %v", err)
		}

		defer gracefullShutDown(server)
	}()

	select {
	case success := <-tokenReady:
		if !success {
			return nil, fmt.Errorf("Problem of the authentication, please try again")
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("Timeout reached, forcing shutdown")
	}

	return token, nil
}

func configSpotifyAuthServer(tokenReady chan bool, token *entities.SpotifyToken) (*http.Server, error) {
	// Start the auth server in a goroutine
	app, err := config.NewApplication()
	if err != nil {
		return nil, fmt.Errorf("Error creating application: %w", err)
	}

	depend, err := dependencies.BuildDependencies()
	if err != nil {
		return nil, fmt.Errorf("Error building dependencies: %w", err)
	}

	// Override the callback handler to capture token and signal completion
	authHandler := depend.SpotifyHandler.Auth.(*handler.SpotifyAuthHandler)
	authHandler.CallBackFunc = func(w http.ResponseWriter, r *http.Request) {

		if authHandler.Auth.Token.AccessToken == "" {
			tokenReady <- false
			return
		}

		*token = authHandler.Auth.Token

		tokenReady <- true
	}

	MapRoutes(app.Server.Mux, depend)

	return &http.Server{
		Addr:    ":8080",
		Handler: app.Server.Mux,
	}, nil
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
	fmt.Printf("📱 If it doesn't open automatically, please navigate to: %s  \n", url)

	// URL is defined in the printed message above
	fmt.Printf("\n⏳ Waiting for authorization...\n\n\n")
}

func gracefullShutDown(server *http.Server) {
	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		server.Close()
	}
}
