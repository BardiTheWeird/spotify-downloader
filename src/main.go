package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"spotify-downloader/src/odeslii"
	"spotify-downloader/src/spotify"
	"strings"
)

var (
	appConfig AppConfig
)

type AppConfig struct {
	ClientId     string
	ClientSecret string
}

func B64Strict(s string) string {
	return base64.RawStdEncoding.Strict().EncodeToString([]byte(s))
}

func (config AppConfig) GetB64() string {
	return B64Strict(config.ClientId + ":" + config.ClientSecret)
}

func getEnvOrDefault(envKey, def string) string {
	val, err := os.LookupEnv(envKey)
	if !err {
		fmt.Printf("%s is not set, using %s\n", envKey, def)
		val = def
	}
	return val
}

func configureApp() {
	appConfig = AppConfig{
		ClientId:     getEnvOrDefault("CLIENT_ID", "00000000000000000000000000000000"),
		ClientSecret: getEnvOrDefault("CLIENT_SECRET", "00000000000000000000000000000000"),
	}
}

func main() {
	configureApp()
	spotify.Authenticate(appConfig.GetB64())

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("there might be a list of all the endpoints here sometime in the future"))
	})

	http.HandleFunc("/playlist/", func(rw http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/playlist/"):]

		playlist, statusCode, err := spotify.GetPlaylistById(id)
		if statusCode == 401 {
			spotify.Authenticate(appConfig.GetB64())
			playlist, statusCode, err = spotify.GetPlaylistById(id)
		}
		if statusCode != 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		rw.Write([]byte("Tracks:\n------\n"))
		for _, v := range playlist.Tracks.Items {
			track := v.Track
			artistStrings := make([]string, 0, len(track.Artists))
			for _, artist := range track.Artists {
				artistStrings = append(artistStrings, artist.Name)
			}

			rw.Write([]byte(fmt.Sprintf(
				"Title: %s\n"+
					"Artist: %s\n"+
					"Image: %s\n"+
					"URL: %s\n"+
					"Id: %s\n",
				track.Name,
				strings.Join(artistStrings, "; "),
				track.Album.Images[0].Url,
				track.Href,
				track.Id)))
			rw.Write([]byte("---\n"))
		}
	})

	http.HandleFunc("/s2y/", func(rw http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/s2y/"):]

		link, exists, err := odeslii.GetYoutubeLinkBySpotifyId(id)
		if err != nil {
			log.Default().Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			rw.Write([]byte(fmt.Sprintf("Odeslii can't find a YouTube link for a track with a Spotify id %s\n", link)))
			return
		}

		rw.Write([]byte(fmt.Sprintf("Here's your link: %s", link)))
	})

	log.Println("Starting a server at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
