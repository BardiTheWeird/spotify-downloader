package downloader

import (
	"spotify-downloader/models"
	"sync"
)

type downloadEntrySyncMap struct {
	downloadEntries sync.Map
}

func (m *downloadEntrySyncMap) Store(trackId string, value models.DownloadEntry) {
	m.downloadEntries.Store(trackId, value)
}

func (m *downloadEntrySyncMap) Load(trackId string) (models.DownloadEntry, bool) {
	val, ok := m.downloadEntries.Load(trackId)
	if !ok {
		return models.DownloadEntry{}, ok
	}
	return val.(models.DownloadEntry), ok
}
