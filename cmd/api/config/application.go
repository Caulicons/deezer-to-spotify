package config

import (
	"fmt"
	"net/http"
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
	http.ListenAndServe(a.Server.address, a.Server.Mux)
}
