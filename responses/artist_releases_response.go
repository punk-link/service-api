package responses

type ArtistReleasesResponse struct {
	ArtistSpotifyId string `json:"artistSpotifyId"`
	Items           []ArtistRelease
}
