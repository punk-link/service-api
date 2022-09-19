package artists

type SlimTrack struct {
	ArtistNames string `json:"artistNames"`
	IsExplicit  bool   `json:"explicit"`
	Name        string `json:"name"`
}
