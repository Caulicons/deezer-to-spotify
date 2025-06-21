package config

import (
	"fmt"
)

type Application struct {
	Server *Server
}

func NewApplication() (*Application, error) {

	server := NewServerConfig()

	return &Application{
		Server: server,
	}, nil
}

func (a *Application) Run() {

	fmt.Printf("Server listening on %s - http://127.0.0.1%s/spotify/auth\n", a.Server.address, a.Server.address)

	// srv := &http.Server{Addr: ":8080"}

	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Server failed: %v", err)
	// 	}
	// }()

	// // Wait for a signal (like Ctrl+C)
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// <-quit

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatalf("Server Shutdown Failed:%+v", err)
	// }
	// log.Println("Server exited properly")
}
