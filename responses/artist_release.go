package responses

type ArtistRelease struct {
	SpotifyId     string          `json:"spotifyId"`
	ImageMetadata []ImageMetadata `json:"images"`
	Name          string          `json:"name"`
	ReleaseDate   string          `json:"releaseDate"`
	TrackNumber   int             `json:"trackNumber"`
	Type          string          `json:"type"`
}
