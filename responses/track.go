package responses

type Track struct {
	SpotifyId       string               `json:"id"`
	Artists         []ArtistSearchResult `json:"artists"`
	DiscNumber      int                  `json:"discNumber"`
	DurationSeconds int                  `json:"durationSeconds"`
	IsExplicit      bool                 `json:"explicit"`
	Name            string               `json:"name"`
	TrackNumber     int                  `json:"trackNumber"`
}
