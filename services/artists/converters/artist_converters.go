package converters

import (
	data "main/data/artists"
	"main/helpers"
	models "main/models/artists"
	spotifyArtists "main/models/spotify/artists"
	spotifyReleases "main/models/spotify/releases"
	commonConverters "main/services/common/converters"
	"time"
)

func ToArtist(dbArtist data.Artist) (models.Artist, error) {
	imageDetails, err := commonConverters.FromJson(dbArtist.ImageDetails)

	// TODO
	//releases, releaseErr := ToReleases(dbArtist.Releases)
	//err = helpers.CombineErrors(err, releaseErr)

	return models.Artist{
		Id:           dbArtist.Id,
		ImageDetails: imageDetails,
		LabelId:      dbArtist.LabelId,
		Name:         dbArtist.Name,
		//Releases:     releases,
	}, err
}

func ToArtists(dbArtists []data.Artist) ([]models.Artist, error) {
	var err error
	results := make([]models.Artist, len(dbArtists))
	for i, dbArtist := range dbArtists {
		artist, localErr := ToArtist(dbArtist)
		if localErr != nil {
			err = helpers.CombineErrors(err, localErr)
		}

		results[i] = artist
	}

	return results, err
}

func ToArtistsFromSpotifyTrack(track spotifyReleases.Track, artists map[string]data.Artist) []models.Artist {
	trackArtists := make([]models.Artist, len(track.Artists))
	for i, artist := range track.Artists {
		trackArtists[i] = models.Artist{
			Id:           artists[artist.Id].Id,
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
			Name:         artist.Name,
		}
	}

	return trackArtists
}

func ToArtistSearchResults(spotifyArtists []spotifyArtists.SlimArtist) []models.ArtistSearchResult {
	results := make([]models.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = models.ArtistSearchResult{
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func ToDbArtist(artist spotifyArtists.Artist, labelId int, timeStamp time.Time) (data.Artist, error) {
	imageDetailsJson, err := commonConverters.ToJsonFromSpotify(artist.ImageDetails, artist.Name)

	return data.Artist{
		Created:      timeStamp,
		ImageDetails: imageDetailsJson,
		LabelId:      labelId,
		Name:         artist.Name,
		SpotifyId:    artist.Id,
		Updated:      timeStamp,
	}, err
}
