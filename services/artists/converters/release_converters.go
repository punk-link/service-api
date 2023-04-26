package converters

import (
	"encoding/json"
	"fmt"
	commonConverters "main/converters/common"
	artistData "main/data/artists"
	"main/helpers"
	artistModels "main/models/artists"
	"main/models/artists/enums"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"time"
)

func ToDbReleaseFromSpotify(release releaseSpotifyPlatformModels.Release, artists map[string]artistData.Artist, timeStamp time.Time) (artistData.Release, error) {
	imageDetailsJson, err := commonConverters.ToJsonFromSpotify(release.ImageDetails, release.Name)

	tracks := ToTracksFromSpotify(release.Tracks.Items, artists)
	tracksJson, err := getTracksJson(err, tracks)
	releaseDate, err := getReleaseDate(err, release)
	releaseArtistIdsJson, err := getReleaseArtistIdsJson(err, release, artists)
	featuringArtistIdsJson, err := getFeaturingArtistIdsJson(err, release, artists)
	releaseType, err := getReleaseType(err, release.Type)
	if err != nil {
		return artistData.Release{}, err
	}

	return artistData.Release{
		Created:            timeStamp,
		Description:        "",
		FeaturingArtistIds: featuringArtistIdsJson,
		ImageDetails:       imageDetailsJson,
		Label:              getLabelName(release),
		Name:               release.Name,
		ReleaseArtistIds:   releaseArtistIdsJson,
		ReleaseDate:        releaseDate,
		SpotifyId:          release.Id,
		TrackNumber:        release.TrackNumber,
		Tracks:             tracksJson,
		Type:               releaseType,
		Upc:                release.ExternalIds.Upc,
		Updated:            timeStamp,
	}, nil
}

func ToReleases(dbReleases []artistData.Release, artists map[int]artistModels.Artist, tags map[int][]artistModels.Tag) ([]artistModels.Release, error) {
	var err error

	results := make([]artistModels.Release, len(dbReleases))
	for i, dbRelease := range dbReleases {
		featuringArtists, featuringArtistErr := ToArtistsFromIds(dbRelease.FeaturingArtistIds, artists)
		releaseArtists, releaseArtistErr := ToArtistsFromIds(dbRelease.ReleaseArtistIds, artists)
		tracks, tracksErr := getTracks(dbRelease.Tracks)
		imageDetails, imageErr := commonConverters.FromJson(dbRelease.ImageDetails)

		err = helpers.CombineErrors(err, helpers.AccumulateErrors(featuringArtistErr, releaseArtistErr, tracksErr, imageErr))

		releaseTags := getReleaseTags(tags, dbRelease.Id)

		results[i] = artistModels.Release{
			Id:               dbRelease.Id,
			Description:      dbRelease.Description,
			FeaturingArtists: featuringArtists,
			ImageDetails:     imageDetails,
			Label:            dbRelease.Label,
			Name:             dbRelease.Name,
			ReleaseArtists:   releaseArtists,
			ReleaseDate:      dbRelease.ReleaseDate,
			Tags:             releaseTags,
			TrackNumber:      dbRelease.TrackNumber,
			Tracks:           tracks,
			Type:             dbRelease.Type,
		}
	}

	return results, err
}

func getFeaturingArtistIds(release releaseSpotifyPlatformModels.Release, artists map[string]artistData.Artist) []int {
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

func getFeaturingArtistIdsJson(err error, release releaseSpotifyPlatformModels.Release, artists map[string]artistData.Artist) (string, error) {
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

func getLabelName(release releaseSpotifyPlatformModels.Release) string {
	if release.Label == "" {
		return release.Artists[0].Name
	}

	return release.Label
}

func getReleaseArtistIds(release releaseSpotifyPlatformModels.Release, artists map[string]artistData.Artist) []int {
	results := make([]int, len(release.Artists))
	for i, artist := range release.Artists {
		results[i] = artists[artist.Id].Id
	}

	return results
}

func getReleaseArtistIdsJson(err error, release releaseSpotifyPlatformModels.Release, artists map[string]artistData.Artist) (string, error) {
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

func getReleaseDate(err error, release releaseSpotifyPlatformModels.Release) (time.Time, error) {
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

func getReleaseTags(tags map[int][]artistModels.Tag, dbReleaseId int) []artistModels.Tag {
	if result, isExist := tags[dbReleaseId]; isExist {
		return result
	}

	return make([]artistModels.Tag, 0)
}

func getReleaseType(err error, target string) (string, error) {
	if err != nil {
		return "", err
	}

	var result string
	switch target {
	case enums.Album:
		result = enums.Album
	case enums.Compilation:
		result = enums.Compilation
	case enums.Single:
		result = enums.Single
	default:
		result = enums.Unknown
	}

	return result, err
}

func getTracks(tracksJson string) ([]artistModels.Track, error) {
	var tracks []artistModels.Track
	err := json.Unmarshal([]byte(tracksJson), &tracks)
	if err != nil {
		tracks = make([]artistModels.Track, 0)
	}

	return tracks, err
}

func getTracksJson(err error, tracks []artistModels.Track) (string, error) {
	if err != nil {
		return "", err
	}

	tracksBytes, err := json.Marshal(tracks)
	if err != nil {
		return "", fmt.Errorf("can't serialize track: '%s'", err.Error())
	}

	return string(tracksBytes), nil
}
