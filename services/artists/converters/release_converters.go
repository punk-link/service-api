package converters

import (
	"encoding/json"
	"fmt"
	data "main/data/artists"
	"main/helpers"
	models "main/models/artists"
	spotifyReleases "main/models/spotify/releases"
	commonConverters "main/services/common/converters"
	"time"
)

func ToReleases(dbReleases []data.Release) ([]models.Release, error) {
	var err error

	results := make([]models.Release, len(dbReleases))
	for i, dbRelease := range dbReleases {
		featuringArtists, featuringArtistErr := ToArtists(dbRelease.Artists)
		releaseArtists, releaseArtistErr := ToArtists(dbRelease.ReleaseArtists)
		tracks, tracksErr := getTracks(dbRelease.Tracks)
		imageDetails, imageErr := commonConverters.FromJson(dbRelease.ImageDetails)

		err = helpers.CombineErrors(err, helpers.AccumulateErrors(featuringArtistErr, releaseArtistErr, tracksErr, imageErr))

		results[i] = models.Release{
			Id:               dbRelease.Id,
			FeaturingArtists: featuringArtists,
			ImageDetails:     imageDetails,
			Lable:            dbRelease.Label,
			Name:             dbRelease.Name,
			ReleaseArtists:   releaseArtists,
			ReleaseDate:      dbRelease.ReleaseDate,
			TrackNumber:      dbRelease.TrackNumber,
			Tracks:           tracks,
			Type:             dbRelease.Type,
		}
	}

	return results, err
}

func ToDbReleaseFromSpotify(release spotifyReleases.Release, artists map[string]data.Artist, timeStamp time.Time) (data.Release, error) {
	imageDetailsJson, err := commonConverters.ToJsonFromSpotify(release.ImageDetails, release.Name)

	tracks := ToTracksFromSpotify(release.Tracks.Items, artists)
	tracksJson, err := getTracksJson(err, tracks)
	releaseDate, err := getReleaseDate(err, release)
	releaseArtistIdsJson, err := getReleaseArtistIdsJson(err, release, artists)
	featuringArtistIdsJson, err := getFeaturingArtistIdsJson(err, release, artists)
	if err != nil {
		return data.Release{}, err
	}

	return data.Release{
		Created:            timeStamp,
		FeaturingArtistIds: featuringArtistIdsJson,
		ImageDetails:       imageDetailsJson,
		Label:              getLabelName(release),
		Name:               release.Name,
		ReleaseArtistIds:   releaseArtistIdsJson,
		ReleaseDate:        releaseDate,
		SpotifyId:          release.Id,
		TrackNumber:        release.TrackNumber,
		Tracks:             tracksJson,
		Type:               release.Type,
		Updated:            timeStamp,
	}, nil
}

func getFeaturingArtistIds(release spotifyReleases.Release, artists map[string]data.Artist) []int {
	featuredArtistIds := make(map[string]int)
	for _, artist := range release.Artists {
		if _, isExists := featuredArtistIds[artist.Id]; isExists {
			continue
		}

		featuredArtistIds[artist.Id] = artists[artist.Id].Id
	}

	for _, track := range release.Tracks.Items {
		for _, artist := range track.Artists {
			if _, isExists := featuredArtistIds[artist.Id]; isExists {
				continue
			}

			featuredArtistIds[artist.Id] = artists[artist.Id].Id
		}
	}

	results := make([]int, 0)
	for _, id := range featuredArtistIds {
		results = append(results, id)
	}

	return results
}

func getFeaturingArtistIdsJson(err error, release spotifyReleases.Release, artists map[string]data.Artist) (string, error) {
	if err != nil {
		return "", err
	}

	featuringArtistIds := getFeaturingArtistIds(release, artists)
	bytes, err := json.Marshal(featuringArtistIds)
	if err != nil {
		return "", fmt.Errorf("can't serialize featuring artist ids: '%s'", err.Error())
	}

	return string(bytes), nil
}

func getLabelName(release spotifyReleases.Release) string {
	if release.Label == "" {
		return release.Artists[0].Name
	}

	return release.Label
}

func getReleaseArtistIds(release spotifyReleases.Release, artists map[string]data.Artist) []int {
	results := make([]int, len(release.Artists))
	for i, artist := range release.Artists {
		results[i] = artists[artist.Id].Id
	}

	return results
}

func getReleaseArtistIdsJson(err error, release spotifyReleases.Release, artists map[string]data.Artist) (string, error) {
	if err != nil {
		return "", err
	}

	releaseArtistIds := getReleaseArtistIds(release, artists)
	bytes, err := json.Marshal(releaseArtistIds)
	if err != nil {
		return "", fmt.Errorf("can't serialize release artist ids: '%s'", err.Error())
	}

	return string(bytes), nil
}

func getReleaseDate(err error, release spotifyReleases.Release) (time.Time, error) {
	if err != nil {
		return time.Time{}, err
	}

	format := time.RFC3339
	switch release.ReleaseDatePrecision {
	case "day":
		format = "2006-01-02"
	case "month":
		format = "2006-01"
	case "year":
		format = "2006"
	}

	releaseDate, err := time.Parse(format, release.ReleaseDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("spotify date format parsing error: '%s'", err.Error())
	}

	return releaseDate, nil
}

func getTracks(tracksJson string) ([]models.Track, error) {
	var tracks []models.Track
	err := json.Unmarshal([]byte(tracksJson), &tracks)
	if err != nil {
		tracks = make([]models.Track, 0)
	}

	return tracks, err
}

func getTracksJson(err error, tracks []models.Track) (string, error) {
	if err != nil {
		return "", err
	}

	tracksBytes, err := json.Marshal(tracks)
	if err != nil {
		return "", fmt.Errorf("can't serialize track: '%s'", err.Error())
	}

	return string(tracksBytes), nil
}
