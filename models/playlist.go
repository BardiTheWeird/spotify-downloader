package models

import "spotify-downloader/spotify"

type Track struct {
	Id      string   `json:"id"`
	Title   string   `json:"title"`
	Artists []string `json:"artists"`

	AlbumTitle string `json:"album_title"`
	AlbumImage string `json:"album_image"`
	AlbumHref  string `json:"album_href"`

	Href       string `json:"href"`
	PreviewUrl string `json:"preview_url"`
}

type PlaylistOwner struct {
	Id          string `json:"id"`
	DisplayName string `json:"display_name"`
	Href        string `json:"href"`
}

type Playlist struct {
	Id     string        `json:"id"`
	Name   string        `json:"name"`
	Owner  PlaylistOwner `json:"owner"`
	Href   string        `json:"href"`
	Image  string        `json:"image"`
	Tracks []Track       `json:"tracks"`
}

func FromSpotifyPlaylist(p spotify.Playlist) Playlist {
	owner := PlaylistOwner{
		Id:          p.Owner.Id,
		DisplayName: p.Owner.Display_name,
		Href:        p.Owner.Href,
	}

	tracks := make([]Track, 0, len(p.Tracks.Items))
	for _, v := range p.Tracks.Items {
		t := v.Track
		artists := make([]string, 0, len(t.Artists))
		for _, v := range t.Artists {
			artists = append(artists, v.Name)
		}
		albumImage := ""
		if len(t.Album.Images) > 0 {
			albumImage = t.Album.Images[0].Url
		}
		tracks = append(tracks, Track{
			Id:      t.Id,
			Title:   t.Name,
			Artists: artists,

			AlbumTitle: t.Album.Name,
			AlbumImage: albumImage,
			AlbumHref:  t.Album.Href,

			Href:       t.Href,
			PreviewUrl: t.Preview_url,
		})
	}

	playlistImage := ""
	if len(p.Images) > 0 {
		playlistImage = p.Images[0].Url
	}
	return Playlist{
		Id:     p.Id,
		Name:   p.Name,
		Owner:  owner,
		Href:   p.Href,
		Image:  playlistImage,
		Tracks: tracks,
	}
}
