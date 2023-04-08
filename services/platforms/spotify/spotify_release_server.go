package spotify

import releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"

type SpotifyReleaseServer interface {
	Get(spotifyIds []string) []releaseSpotifyPlatformModels.Release
	GetByArtistId(spotifyId string) []releaseSpotifyPlatformModels.Release
}
