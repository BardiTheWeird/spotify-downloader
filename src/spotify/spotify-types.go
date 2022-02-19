package spotify

import "fmt"

type ClientToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Followers_Type struct {
	Total int
}

type Image struct {
	Url string
}

type Owner_Type struct {
	Display_name  string
	External_urls map[string]string
	Href          string
	Id            string
	Type          string
	Uri           string
}

type AddedBy_Type struct {
	External_urls map[string]string
	Href          string
	Id            string
	Type          string
	Uri           string
}

type Artist struct {
	External_urls map[string]string
	Href          string
	Id            string
	Name          string
	Type          string
	Uri           string
}

type Album_Type_Go struct {
	Artists []Artist
	Href    string
	Id      string
	Images  []Image
	Name    string
}

type Track_Type struct {
	Album        Album_Type_Go
	Artists      []Artist
	Duration_ms  int
	Href         string
	Id           string
	Name         string
	Preview_url  string
	Track_number int
}

func (t Track_Type) Print() {
	fmt.Print("Artists: ")
	for _, v := range t.Artists {
		fmt.Print(v.Name)
	}
	fmt.Println()

	fmt.Printf("Id: %s\n", t.Id)
	fmt.Printf("Album: %s\n", t.Album.Name)
	fmt.Printf("Duration (ms): %d\n", t.Duration_ms)
	fmt.Printf("Href: %s\n", t.Href)
	fmt.Printf("Track_number: %d\n", t.Track_number)
}

type TrackEntry struct {
	Added_at        string
	Added_by        AddedBy_Type
	Is_local        bool
	Track           Track_Type
	Video_thumbnail map[string]string
}

func (t TrackEntry) Print() {
	t.Track.Print()
}

type Tracks_Type struct {
	Href     string
	Items    []TrackEntry
	Limit    int
	Next     string
	Offset   int
	Previous string
	Total    int
}

type Playlist struct {
	Id          string
	Name        string
	Description string
	Type        string
	Owner       Owner_Type
	Href        string
	Images      []Image
	Tracks      Tracks_Type
}

func (p Playlist) Print() {
	fmt.Println("_Playlist:")
	fmt.Printf("Id: %s\n", p.Id)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Owner: %s\n", p.Owner.Display_name)

	fmt.Printf("Href: %s\n", p.Href)
	fmt.Printf("Images: %v\n", p.Images)

	fmt.Println("__Tracks:")
	for _, track := range p.Tracks.Items {
		track.Print()
		fmt.Println("---")
	}
}
