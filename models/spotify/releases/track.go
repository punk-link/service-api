package releases

import "main/models/spotify/search"

type Track struct {
	Id                   string          `json:"id"`
	Artists              []search.Artist `json:"artists"`
	DiscNumber           int             `json:"disc_number"`
	DurationMilliseconds int             `json:"duration_ms"`
	IsExplicit           bool            `json:"explicit"`
	Name                 string          `json:"name"`
	TrackNumber          int             `json:"track_number"`
}
