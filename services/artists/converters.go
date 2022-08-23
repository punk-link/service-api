package artists

import (
	data "main/data/artists"
	models "main/models/artists"
	"main/models/common"
	spotifyArtists "main/models/spotify/artists"
	imageConverter "main/services/common/converters"
	"time"
)

func toArtist(dbArtist data.Artist) models.Artist {

	return models.Artist{
		Id:           dbArtist.Id,
		ImageDetails: common.ImageDetails{},
		LabelId:      dbArtist.LabelId,
		Name:         dbArtist.Name,
		Releases:     []models.Release{},
	}
}

func toArtistSearchResults(spotifyArtists []spotifyArtists.SlimArtist) []models.ArtistSearchResult {
	results := make([]models.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = models.ArtistSearchResult{
			ImageDetails: imageConverter.ToImageDetailsFromSpotify(artist.ImageDetails),
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func toDbArtist(artist spotifyArtists.Artist, labelId int, timeStamp time.Time) data.Artist {
	return data.Artist{
		Created:   timeStamp,
		LabelId:   labelId,
		Name:      artist.Name,
		SpotifyId: artist.Id,
		Updated:   timeStamp,
	}
}
