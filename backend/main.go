package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"spotify-downloader/server"
)

func main() {
	runServer()
}

func runServer() {
	srv := server.Server{}
	srv.ConfigureDefaults()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening on port",
		listener.Addr().(*net.TCPAddr).Port)

	log.Fatal(http.Serve(listener, &srv))
}
