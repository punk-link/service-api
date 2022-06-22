package responses

type ArtistSearchResponse struct {
	Items []ArtistSearchResult `json:"items"`
}

type ArtistSearchResult struct {
	SpotifyId string `json:"spotifyId"`
	Name      string `json:"name"`
}
