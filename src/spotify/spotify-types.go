package spotify

import "fmt"

type ClientToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Followers_Type struct {
	// Href  string
	Total int
}

type Image struct {
	// Height int
	Url string
	// Width  int
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
	// Album_type             string
	Artists []Artist
	// Available_markets      []string
	// External_urls          map[string]string
	Href   string
	Id     string
	Images []Image
	Name   string
	// Release_date           string
	// Release_date_precision string
	// Total_tracks           int
	// Type                   string
	// Uri                    string
}

type Track_Type struct {
	Album   Album_Type_Go
	Artists []Artist
	// Available_markets []string
	// Disc_number       int
	Duration_ms int
	// Episode           bool
	// Explicit          bool
	// External_ids      map[string]string
	// External_urls     map[string]string
	Href string
	Id   string
	// Is_local          bool
	Name string
	// Popularity        int
	Preview_url string
	// Track             bool
	Track_number int
	// Type              string
	// Uri               string
}

func (t Track_Type) Print() {
	fmt.Print("Artists: ")
	for _, v := range t.Artists {
		fmt.Print(v.Name)
	}
	fmt.Println()

	fmt.Printf("Album: %s\n", t.Album.Name)
	fmt.Printf("Duration (ms): %d\n", t.Duration_ms)
	fmt.Printf("Href: %s\n", t.Href)
	fmt.Printf("Track_number: %d\n", t.Track_number)
}

type TrackEntry struct {
	Added_at string
	Added_by AddedBy_Type
	Is_local bool
	// "primary_color" : null,
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

	// Public        bool
	// Collaborative bool

	// Followers Followers_Type
	Href   string
	Images []Image
	// "primary_color" : null,
	// Snapshot_id   string
	// Uri           string
	// External_urls map[string]string

	Tracks Tracks_Type
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
