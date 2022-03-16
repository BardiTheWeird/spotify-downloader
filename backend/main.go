package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening on port",
		listener.Addr().(*net.TCPAddr).Port)

	log.Fatal(http.Serve(listener, &srv))
}
