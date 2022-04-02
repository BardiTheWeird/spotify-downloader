package publicauthenticator

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	ClientId     string
	ClientSecret string
}

func GetConfigurationFromEnv() Configuration {
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		log.Fatal("CLIENT_ID was not provided")
	}
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("CLIENT_SECRET was not provided")
	}

	return Configuration{clientId, clientSecret}
}

func SpotifyClientAuth(rw http.ResponseWriter, r *http.Request) {
	config := GetConfigurationFromEnv()
	b64Config := base64.RawStdEncoding.Strict().EncodeToString(
		[]byte(config.ClientId + ":" + config.ClientSecret))

	url := "https://accounts.spotify.com/api/token?grant_type=client_credentials"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+b64Config)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	rw.WriteHeader(res.StatusCode)
	rw.Write(body)
}
