package spotify

import artistSpotifyPlatformModels "main/models/platforms/spotify/artists"

type SpotifyArtistServer interface {
	Get(spotifyIds []string) []artistSpotifyPlatformModels.Artist
	GetOne(spotifyId string) (*artistSpotifyPlatformModels.Artist, error)
	Search(query string) []artistSpotifyPlatformModels.SlimArtist
}
