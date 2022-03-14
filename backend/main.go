package main

import (
	"flag"
	"log"
	"net/http"
	"spotify-downloader/server"
)

func main() {
	runServer()
}

func getSettingsPath() string {
	settingsPath := flag.String("settings", "settings.json", "location of a settings.json file")
	flag.Parse()
	log.Println("settings.json is at", *settingsPath)
	return *settingsPath
}

func runServer() {
	srv := server.Server{}
	srv.SettingsFileLocation = getSettingsPath()
	srv.ConfigureFromSettingsFile()
	srv.ConfigureRoutes()
	srv.DiscoverFeatures()

	srv.SonglinkHelper.SetDefaults()

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", &srv))
}
