package downloader

import (
	"spotify-downloader/models"
	"sync"
)

type downloadEntrySyncMap struct {
	downloadEntries sync.Map
}

func (m *downloadEntrySyncMap) Store(key string, value models.DownloadEntry) {
	m.downloadEntries.Store(key, value)
}

func (m *downloadEntrySyncMap) Load(key string) (models.DownloadEntry, bool) {
	val, ok := m.downloadEntries.Load(key)
	if !ok {
		return models.DownloadEntry{}, ok
	}
	return val.(models.DownloadEntry), ok
}
