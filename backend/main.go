package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"spotify-downloader/server"
	"strconv"
)

func main() {
	runServer()
}

func runServer() {
	srv := server.Server{}
	flag.StringVar(&srv.SettingsPath, "settings-path", "", "path to settings.json")
	flag.StringVar(&srv.SpotifyHelper.PublicAuthorizationEndpoint,
		"authorization-endpoint", "",
		"url to send a GET request to if client does not provide OAuth credentials")
	var port string
	flag.StringVar(&port, "port", "0", "port to listen at")
	flag.Parse()
	srv.ConfigureDefaults()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	if port == "0" {
		port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	}
	fmt.Println("listening on port", port)

	log.Fatal(http.Serve(listener, &srv))
}
