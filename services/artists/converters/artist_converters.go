package converters

import (
	data "main/data/artists"
	models "main/models/artists"
	"main/models/common"
	spotifyArtists "main/models/spotify/artists"
	spotifyReleases "main/models/spotify/releases"
	commonConverters "main/services/common/converters"
	"time"
)

func ToArtist(dbArtist data.Artist) models.Artist {
	return models.Artist{
		Id:           dbArtist.Id,
		ImageDetails: common.ImageDetails{},
		LabelId:      dbArtist.LabelId,
		Name:         dbArtist.Name,
		Releases:     []models.Release{},
	}
}

func ToArtistsFromSpotifyTrack(track spotifyReleases.Track, artists map[string]data.Artist) []models.Artist {
	trackArtists := make([]models.Artist, len(track.Artists))
	for i, artist := range track.Artists {
		trackArtists[i] = models.Artist{
			Id:           artists[artist.Id].Id,
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails),
			Name:         artist.Name,
		}
	}

	return trackArtists
}

func ToArtistSearchResults(spotifyArtists []spotifyArtists.SlimArtist) []models.ArtistSearchResult {
	results := make([]models.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = models.ArtistSearchResult{
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails),
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func ToDbArtist(artist spotifyArtists.Artist, labelId int, timeStamp time.Time) data.Artist {
	return data.Artist{
		Created:   timeStamp,
		LabelId:   labelId,
		Name:      artist.Name,
		SpotifyId: artist.Id,
		Updated:   timeStamp,
	}
}
