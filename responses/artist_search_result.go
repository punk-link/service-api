package responses

type ArtistSearchResult struct {
	SpotifyId     string          `json:"spotifyId"`
	ImageMetadata []ImageMetadata `json:"images"`
	Name          string          `json:"name"`
}
