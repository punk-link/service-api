package converters

import (
	"encoding/json"
	commonConverters "main/converters/common"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	artistSpotifyPlatformModels "main/models/platforms/spotify/artists"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"time"
)

func ToArtist(dbArtist artistData.Artist, releases []artistModels.Release) (artistModels.Artist, error) {
	imageDetails, err := commonConverters.FromJson(dbArtist.ImageDetails)

	return artistModels.Artist{
		Id:           dbArtist.Id,
		ImageDetails: imageDetails,
		LabelId:      dbArtist.LabelId,
		Name:         dbArtist.Name,
		Releases:     releases,
	}, err
}

func ToArtistsFromIds(artistIdJson string, artists map[int]artistModels.Artist) ([]artistModels.Artist, error) {
	var artistIds []int
	err := json.Unmarshal([]byte(artistIdJson), &artistIds)
	if err != nil {
		return make([]artistModels.Artist, 0), err
	}

	results := make([]artistModels.Artist, 0)
	for _, id := range artistIds {
		if artist, isExists := artists[id]; isExists {
			results = append(results, artist)
		}
	}

	return results, err
}

func ToArtistsFromSpotifyTrack(track releaseSpotifyPlatformModels.Track, artists map[string]artistData.Artist) []artistModels.Artist {
	trackArtists := make([]artistModels.Artist, len(track.Artists))
	for i, artist := range track.Artists {
		trackArtists[i] = artistModels.Artist{
			Id:           artists[artist.Id].Id,
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
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
			ImageDetails: commonConverters.ToImageDetailsFromSpotify(artist.ImageDetails, artist.Name),
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func ToDbArtist(artist *artistSpotifyPlatformModels.Artist, labelId int, timeStamp time.Time) (artistData.Artist, error) {
	imageDetailsJson, err := commonConverters.ToJsonFromSpotify(artist.ImageDetails, artist.Name)

	return artistData.Artist{
		Created:      timeStamp,
		ImageDetails: imageDetailsJson,
		LabelId:      labelId,
		Name:         artist.Name,
		SpotifyId:    artist.Id,
		Updated:      timeStamp,
	}, err
}
