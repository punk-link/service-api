package artists

type Track struct {
	Artists         []Artist `json:"artists"`
	DiscNumber      int      `json:"discNumber"`
	DurationSeconds int      `json:"durationSeconds"`
	IsExplicit      bool     `json:"explicit"`
	Name            string   `json:"name"`
	SpotifyId       string   `json:"spotifyId"`
	TrackNumber     int      `json:"trackNumber"`
}
