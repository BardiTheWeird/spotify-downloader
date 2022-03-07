package server

import (
	"encoding/json"
	"log"
	"os"
)

type ServerSettings struct {
	SpotifyClientId     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"spotify_client_secret"`
}

func (s *Server) ConfigureFromSettingsFile() bool {
	if _, err := os.Stat("settings.json"); os.IsNotExist(err) {
		log.Println("settings.json does not exist")
		return false
	}

	file, err := os.Open("settings.json")
	if err != nil {
		log.Println("error opening settings.json:", err)
		return false
	}
	defer file.Close()

	var settings ServerSettings
	json.NewDecoder(file).Decode(&settings)

	s.SpotifyHelper.ClientId = settings.SpotifyClientId
	s.SpotifyHelper.ClientSecret = settings.SpotifyClientSecret
	log.Println("read settings from settings.json")
	return true
}

func (s *Server) UpdateSettingsFile() bool {
	file, err := os.Create("settings.json")
	if err != nil {
		log.Println("error creating settings.json:", err)
		return false
	}
	defer file.Close()

	settings := ServerSettings{
		SpotifyClientId:     s.SpotifyHelper.ClientId,
		SpotifyClientSecret: s.SpotifyHelper.ClientSecret,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(settings)
	if err != nil {
		log.Println("error encoding server settings:", err)
	}
	log.Println("saved settings to settings.json")
	return true
}
