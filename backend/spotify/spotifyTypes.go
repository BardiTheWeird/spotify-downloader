package spotify

import (
	"spotify-downloader/models"
)

type albumTracks struct {
	Name   string
	Images []struct {
		Url string
	}
	Tracks struct {
		Items []track
	}
}

func (a *albumTracks) toTracks() []track {
	tracks := a.Tracks.Items
	for i := 0; i < len(tracks); i++ {
		tracks[i].Album.Name = a.Name
		tracks[i].Album.Images = a.Images
	}
	return tracks
}

type track struct {
	Id      string
	Name    string
	Artists []struct {
		Name string
	}

	Album struct {
		Name   string
		Images []struct {
			Url string
		}
	}

	Preview_url string
}

type tracksObject struct {
	Items []struct {
		Track track
	}
}

func (t *tracksObject) toTrackSlice() []track {
	tracks := make([]track, 0, len(t.Items))
	for _, v := range t.Items {
		tracks = append(tracks, v.Track)
	}
	return tracks
}

func toModelsPlaylist(tracksIn []track) models.Playlist {
	tracks := make([]models.Track, 0, len(tracksIn))
	for _, v := range tracksIn {
		t := v
		artists := make([]string, 0, len(t.Artists))
		for _, v := range t.Artists {
			artists = append(artists, v.Name)
		}
		albumImage := ""
		if len(t.Album.Images) > 0 {
			albumImage = t.Album.Images[0].Url
		}
		tracks = append(tracks, models.Track{
			Id:      t.Id,
			Title:   t.Name,
			Artists: artists,

			AlbumTitle: t.Album.Name,
			AlbumImage: albumImage,
			PreviewUrl: t.Preview_url,
		})
	}

	return models.Playlist{
		Tracks: tracks,
	}
}
