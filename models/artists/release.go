package artists

type Release struct {
	SpotifyId     string               `json:"spotifyId"`
	Artists       []ArtistSearchResult `json:"artists"`
	ImageMetadata []ImageMetadata      `json:"images"`
	Lable         string               `json:"label"`
	Name          string               `json:"name"`
	ReleaseDate   string               `json:"releaseDate"`
	TrackNumber   int                  `json:"trackNumber"`
	Tracks        []Track              `json:"tracks"`
	Type          string               `json:"type"`
}
