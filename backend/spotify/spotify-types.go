package spotify

import "spotify-downloader/models"

type playlist struct {
	Id    string
	Name  string
	Href  string
	Owner struct {
		Id           string
		Display_name string
		Href         string
	}
	Images []struct {
		Url string
	}
	Tracks struct {
		Items []struct {
			Added_at string
			Is_local bool
			Track    struct {
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
					Href string
				}

				Href        string
				Preview_url string
			}
		}
	}
}

func (p *playlist) toModelsPlaylist() models.Playlist {
	type Track = struct {
		Id      string   `json:"id"`
		Title   string   `json:"title"`
		Artists []string `json:"artists"`

		AlbumTitle string `json:"album_title"`
		AlbumImage string `json:"album_image"`
		AlbumHref  string `json:"album_href"`

		Href       string `json:"href"`
		PreviewUrl string `json:"preview_url"`
	}

	owner := struct {
		Id          string `json:"id"`
		DisplayName string `json:"display_name"`
		Href        string `json:"href"`
	}{
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
	return models.Playlist{
		Id:     p.Id,
		Name:   p.Name,
		Owner:  owner,
		Href:   p.Href,
		Image:  playlistImage,
		Tracks: tracks,
	}
}
