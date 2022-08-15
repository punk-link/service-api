package artists

import (
	"errors"
	"main/data"
	artistData "main/data/artists"
	"main/models/artists"
	"main/models/labels"
	"main/services/spotify"
	"time"
)

type ArtistService struct {
	spotifyService *spotify.SpotifyService
}

func BuildArtistService(spotifyService *spotify.SpotifyService) *ArtistService {
	return &ArtistService{
		spotifyService: spotifyService,
	}
}

func (t *ArtistService) AddArtist(currentManager labels.ManagerContext, spotifyId string) (interface{}, error) {
	if spotifyId == "" {
		return "", errors.New("artist's spotify ID is empty")
	}

	var dbArtist artistData.Artist
	queryResult := data.DB.Model(&artistData.Artist{}).Preload("Releases").Where("spotify_id = ?", spotifyId).FirstOrInit(&dbArtist)
	if queryResult.Error != nil {
		return labels.Label{}, queryResult.Error
	}

	if dbArtist.Id != 0 {
		if dbArtist.LabelId != currentManager.LabelId {
			return "", errors.New("artist already added to another label")
		}
	}

	dbReleaseIds := make(map[string]int, len(dbArtist.Releases))
	for _, release := range dbArtist.Releases {
		dbReleaseIds[release.SpotifyId] = 0
	}

	spotifyReleases := t.spotifyService.GetArtistReleases(spotifyId)
	missedReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if _, isContains := dbReleaseIds[spotifyRelease.Id]; !isContains {
			missedReleaseSpotifyIds = append(missedReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	releasesDetails := t.spotifyService.GetReleasesDetails(missedReleaseSpotifyIds)
	now := time.Now().UTC()
	dbReleases := make([]artistData.Release, len(releasesDetails))
	for i, details := range releasesDetails {
		releaseDate, _ := time.Parse("2006-01-02", details.ReleaseDate)
		dbRelease := artistData.Release{
			Created:         now,
			Label:           details.Label,
			Name:            details.Name,
			PrimaryArtistId: dbArtist.Id,
			ReleaseDate:     releaseDate,
			SpotifyId:       details.Id,
			TrackNumber:     details.TrackNumber,
			Type:            details.Type,
			Updated:         now,
		}

		dbReleases[i] = dbRelease
	}

	return "", nil
}

func (t *ArtistService) SearchArtist(query string) []artists.ArtistSearchResult {
	var result []artists.ArtistSearchResult

	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return result
	}

	results := t.spotifyService.SearchArtist(query)

	return spotify.ToArtistSearchResults(results)
}
