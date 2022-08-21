package releases

import "main/models/spotify/artists"

type Track struct {
	Id                   string               `json:"id"`
	Artists              []artists.SlimArtist `json:"artists"`
	DiscNumber           int                  `json:"disc_number"`
	DurationMilliseconds int                  `json:"duration_ms"`
	IsExplicit           bool                 `json:"explicit"`
	Name                 string               `json:"name"`
	TrackNumber          int                  `json:"track_number"`
}
