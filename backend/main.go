package main

import (
	"log"
	"net/http"
	"spotify-downloader/server"
)

func main() {
	runServer()
}

func runServer() {
	srv := server.Server{}
	srv.ConfigureFromEnv()
	srv.ConfigureRoutes()

	srv.SonglinkHelper.SetDefaultEndpoint()

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", &srv))
}
