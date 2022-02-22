package models

type Playlist struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Href  string `json:"href"`
	Owner struct {
		Id          string `json:"id"`
		DisplayName string `json:"display_name"`
		Href        string `json:"href"`
	} `json:"owner"`
	Image  string `json:"image"`
	Tracks []struct {
		Id      string   `json:"id"`
		Title   string   `json:"title"`
		Artists []string `json:"artists"`

		AlbumTitle string `json:"album_title"`
		AlbumImage string `json:"album_image"`
		AlbumHref  string `json:"album_href"`

		Href       string `json:"href"`
		PreviewUrl string `json:"preview_url"`
	} `json:"tracks"`
}
