package config

type Application struct {
	Server *Server
}

func NewApplication() (*Application, error) {

	server := NewServerConfig()

	return &Application{
		Server: server,
	}, nil
}
