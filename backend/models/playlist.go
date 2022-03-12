package models

type Playlist struct {
	Tracks []Track `json:"tracks"`
}

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
