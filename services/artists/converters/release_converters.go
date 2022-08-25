package converters

import (
	"encoding/json"
	"fmt"
	data "main/data/artists"
	"main/helpers"
	models "main/models/artists"
	commonModels "main/models/common"
	spotifyReleases "main/models/spotify/releases"
	commonConverters "main/services/common/converters"
	"time"
)

func ToReleases(dbReleases []data.Release) ([]models.Release, error) {
	var err error

	results := make([]models.Release, len(dbReleases))
	for i, dbRelease := range dbReleases {
		artists, artistErr := ToArtists(dbRelease.Artists)
		tracks, tracksErr := getTracks(dbRelease.Tracks)
		imageDetails, imageErr := commonConverters.FromJson(dbRelease.ImageDetails)

		err = helpers.CombineErrors(err, helpers.AccumulateErrors(artistErr, tracksErr, imageErr))

		results[i] = models.Release{
			Id:           dbRelease.Id,
			Artists:      artists,
			ImageDetails: imageDetails,
			Lable:        dbRelease.Label,
			Name:         dbRelease.Name,
			ReleaseDate:  dbRelease.ReleaseDate,
			TrackNumber:  dbRelease.TrackNumber,
			Tracks:       tracks,
			Type:         dbRelease.Type,
		}
	}

	return results, err
}

func ToDbReleaseFromSpotify(release spotifyReleases.Release, releaseArtists []data.Artist, imageDetails commonModels.ImageDetails, tracks []models.Track, timeStamp time.Time) (data.Release, error) {
	imageDetailsJson, err := commonConverters.ToJson(imageDetails)
	tracksJson, err := getTracksJson(err, tracks)
	releaseDate, err := getReleaseDate(err, release)
	if err != nil {
		return data.Release{}, err
	}

	return data.Release{
		Artists:         releaseArtists,
		Created:         timeStamp,
		ImageDetails:    imageDetailsJson,
		Label:           getLabelName(release),
		Name:            release.Name,
		PrimaryArtistId: releaseArtists[0].Id, // TODO: check
		ReleaseDate:     releaseDate,
		SpotifyId:       release.Id,
		TrackNumber:     release.TrackNumber,
		Tracks:          tracksJson,
		Type:            release.Type,
		Updated:         timeStamp,
	}, nil
}

func getLabelName(release spotifyReleases.Release) string {
	if release.Label == "" {
		return release.Artists[0].Name // TODO: check
	}

	return release.Label
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
