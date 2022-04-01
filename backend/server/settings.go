package server

import (
	"encoding/json"
	"os"
)

func (s *Server) ConfigureFromSettingsFile() {
	var settings map[string]string = s.ReadSettingsFile()
	defaultIfNotPresent := func(key, defaultValue string) string {
		val, ok := settings[key]
		if !ok {
			return defaultValue
		}
		return val
	}
	s.Features.Ffmpeg.Path = defaultIfNotPresent("ffmpeg", "ffmpeg")
	s.Features.YoutubeDl.Path = defaultIfNotPresent("youtube_dl", "youtube-dl")
}

func (s *Server) UpdateSettingsFile(f func(*map[string]string)) {
	settings := s.ReadSettingsFile()
	f(&settings)
	s.WriteSettingsFile(settings)
}

func (s *Server) ReadSettingsFile() map[string]string {
	bytes, err := os.ReadFile(s.SettingsPath)
	if err != nil {
		return make(map[string]string)
	}

	var settings map[string]string
	json.Unmarshal(bytes, &settings)
	return settings
}

func (s *Server) WriteSettingsFile(settings map[string]string) {
	file, err := os.Create(s.SettingsPath)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, _ := json.MarshalIndent(settings, "", "  ")
	file.Write(bytes)
}
