package main

import (
	"fmt"

	"github.com/caulicons/deezer-to-spotify/cmd/api/config"
	dependencies "github.com/caulicons/deezer-to-spotify/cmd/api/depedencies"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	err := run()

	if err != nil {
		fmt.Println(err)
	}
}

var baseDeezerURL = "https://api.deezer.com/user/1237568626/tracks?index=0"

func run() error {

	app, err := config.NewApplication()
	if err != nil {
		return err
	}

	depend, err := dependencies.BuildDependencies()
	if err != nil {
		return err
	}

	config.MapRoutes(app.Server.Mux, depend)

	app.Run()
	return nil
}
