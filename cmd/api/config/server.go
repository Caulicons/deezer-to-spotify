package config

import (
	"net/http"
	"os"
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
