package converters

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"
	artistSpotifyPlatformModels "main/models/platforms/spotify/artists"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"main/services/common/converters"
	"time"
)

func ToArtist(dbArtist artistData.Artist, releases []artistModels.Release) (artistModels.Artist, error) {
	imageDetails, err := converters.FromJson(dbArtist.ImageDetails)

	return artistModels.Artist{
		Id:           dbArtist.Id,
		ImageDetails: imageDetails,
		LabelId:      dbArtist.LabelId,
		Name:         dbArtist.Name,
		Releases:     releases,
	}, err
}

func ToArtistsFromSpotifyTrack(track releaseSpotifyPlatformModels.Track, artists map[string]artistData.Artist) []artistModels.Artist {
	trackArtists := make([]artistModels.Artist, len(track.Artists))
	for i, artist := range track.Artists {
		trackArtists[i] = artistModels.Artist{
			Id:           artists[artist.Id].Id,
			ImageDetails: converters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
			Name:         artist.Name,
			Releases:     make([]artistModels.Release, 0),
		}
	}

	return trackArtists
}

func ToArtistSearchResults(spotifyArtists []artistSpotifyPlatformModels.SlimArtist) []artistModels.ArtistSearchResult {
	results := make([]artistModels.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = artistModels.ArtistSearchResult{
			ImageDetails: converters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func ToDbArtist(artist *artistSpotifyPlatformModels.Artist, labelId int, timeStamp time.Time) (artistData.Artist, error) {
	imageDetailsJson, err := converters.ToJsonFromSpotify(artist.ImageDetails, artist.Name)

	return artistData.Artist{
		Created:      timeStamp,
		ImageDetails: imageDetailsJson,
		LabelId:      labelId,
		Name:         artist.Name,
		SpotifyId:    artist.Id,
		Updated:      timeStamp,
	}, err
}
