package main

import (
	"log"
	"net/http"
	"spotify-downloader/server"
	"spotify-downloader/spotify"
)

func main() {
	runServer()
}

func runServer() {
	srv := server.Server{}
	srv.ConfigureFromEnv()
	srv.ConfigureRoutes()

	spotify.Authenticate(srv.GetB64())

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", &srv))
}
