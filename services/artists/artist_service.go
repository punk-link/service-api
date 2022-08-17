package artists

import (
	"errors"
	"main/data"
	artistData "main/data/artists"
	"main/models/artists"
	commonModels "main/models/common"
	"main/models/labels"
	"main/models/spotify/releases"
	"main/services/common"
	"main/services/spotify"
	"sort"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

type ArtistService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
}

func BuildArtistService(logger *common.Logger, spotifyService *spotify.SpotifyService) *ArtistService {
	return &ArtistService{
		logger:         logger,
		spotifyService: spotifyService,
	}
}

func (t *ArtistService) AddArtist(currentManager labels.ManagerContext, spotifyId string) (interface{}, error) {
	if spotifyId == "" {
		return "", errors.New("artist's spotify ID is empty")
	}

	dbArtist, err := t.getDatabaseArtistBySpotifyId(spotifyId)
	if err != nil {
		return "", err
	}

	if dbArtist.Id != 0 {
		if dbArtist.LabelId != currentManager.LabelId {
			return "", errors.New("artist already added to another label")
		}
	}

	_ = t.prepareReleases(dbArtist, spotifyId)

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

func (t *ArtistService) getDatabaseArtistBySpotifyId(spotifyId string) (artistData.Artist, error) {
	var dbArtist artistData.Artist
	queryResult := data.DB.Model(&artistData.Artist{}).Preload("Releases").Where("spotify_id = ?", spotifyId).FirstOrInit(&dbArtist)
	if queryResult.Error != nil {
		return dbArtist, queryResult.Error
	}

	return dbArtist, nil
}

func (t *ArtistService) getImageDetails(release releases.Release) commonModels.ImageDetails {
	details := commonModels.ImageDetails{}
	if 0 < len(release.ImageDetails) {
		sort.SliceStable(release.ImageDetails, func(i, j int) bool {
			return release.ImageDetails[i].Height > release.ImageDetails[j].Height
		})

		details = commonModels.ImageDetails{
			AltText: release.Name,
			Height:  release.ImageDetails[0].Height,
			Url:     release.ImageDetails[0].Url,
			Width:   release.ImageDetails[0].Width,
		}
	}

	return details
}

func (t *ArtistService) getLabelName(release releases.Release) string {
	if release.Label == "" {
		return release.Artists[0].Name
	}

	return release.Label
}

func (t *ArtistService) getMissingReleaseIds(dbArtist artistData.Artist, artistSpotifyId string) []string {
	dbReleaseIds := make(map[string]int, len(dbArtist.Releases))
	for _, release := range dbArtist.Releases {
		dbReleaseIds[release.SpotifyId] = 0
	}

	spotifyReleases := t.spotifyService.GetArtistReleases(artistSpotifyId)
	missingReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if _, isContains := dbReleaseIds[spotifyRelease.Id]; !isContains {
			missingReleaseSpotifyIds = append(missingReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	return missingReleaseSpotifyIds
}

func (t *ArtistService) getReleaseArtists(dbArtist artistData.Artist, release releases.Release) []*artistData.Artist {
	featuredArtists := make([]*artistData.Artist, 0)
	for _, track := range release.Tracks.Items {
		if len(track.Artists) < 2 {
			continue
		}

		trackArtists := make([]*artistData.Artist, len(track.Artists)-1)
		for i := 1; i < len(track.Artists); i++ {
			trackArtists[i-1] = &artistData.Artist{
				Name:      track.Artists[i].Name,
				SpotifyId: track.Artists[i].Id,
			}
		}

		featuredArtists = append(featuredArtists, trackArtists...)
	}

	artists := make([]*artistData.Artist, 1+len(featuredArtists))
	artists[0] = &dbArtist
	artists = append(artists, featuredArtists...)

	return artists
}

func (t *ArtistService) getReleaseDate(release releases.Release) time.Time {
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
		t.logger.LogError(err, "Spotify date format parsing error: '%s'", err.Error())
		return time.Time{}
	}

	return releaseDate
}

func (t *ArtistService) getTracks(release releases.Release) []artists.Track {
	sort.SliceStable(release.Tracks.Items, func(i, j int) bool {
		return release.Tracks.Items[i].DiscNumber < release.Tracks.Items[j].DiscNumber && release.Tracks.Items[i].TrackNumber < release.Tracks.Items[j].TrackNumber
	})

	tracks := make([]artists.Track, len(release.Tracks.Items))
	for i, track := range release.Tracks.Items {
		trackArtists := make([]artists.Artist, len(track.Artists))
		for i, artist := range track.Artists {
			sort.SliceStable(artist.ImageDetails, func(i, j int) bool {
				return artist.ImageDetails[i].Height > artist.ImageDetails[j].Height
			})

			trackArtists[i] = artists.Artist{
				Id:           0,
				ImageDetails: commonModels.ImageDetails{},
				Name:         artist.Name,
			}
		}

		tracks[i] = artists.Track{
			Artists:         trackArtists,
			DiscNumber:      track.DiscNumber,
			DurationSeconds: track.DurationMilliseconds / 1000,
			IsExplicit:      track.IsExplicit,
			Name:            track.Name,
			SpotifyId:       track.Id,
			TrackNumber:     track.TrackNumber,
		}
	}

	return tracks
}

func (t *ArtistService) prepareRelease(timeStamp time.Time, dbArtist artistData.Artist, release releases.Release, wg *sync.WaitGroup, results chan<- artistData.Release) {
	defer wg.Done()

	imageDetails := t.getImageDetails(release)
	imageDetailsJson, err := json.Marshal(imageDetails)
	if err != nil {
		t.logger.LogError(err, "Can't serialize image details: '%s'", err.Error())
		return
	}

	tracks := t.getTracks(release)
	tracksJson, err := json.Marshal(tracks)
	if err != nil {
		t.logger.LogError(err, "Can't serialize track: '%s'", err.Error())
		return
	}

	dbRelease := artistData.Release{
		Artists:         t.getReleaseArtists(dbArtist, release),
		Created:         timeStamp,
		Label:           t.getLabelName(release),
		Name:            release.Name,
		PrimaryArtistId: dbArtist.Id,
		ReleaseDate:     t.getReleaseDate(release),
		SpotifyId:       release.Id,
		TrackNumber:     release.TrackNumber,
		Type:            release.Type,
		Updated:         timeStamp,
	}

	dbRelease.ImageDetails.Set(imageDetailsJson)
	dbRelease.Tracks.Set(tracksJson)

	results <- dbRelease
}

func (t *ArtistService) prepareReleases(dbArtist artistData.Artist, artistSpotifyId string) []artistData.Release {
	missingReleaseSpotifyIds := t.getMissingReleaseIds(dbArtist, artistSpotifyId)
	releasesDetails := t.spotifyService.GetReleasesDetails(missingReleaseSpotifyIds)

	var wg sync.WaitGroup
	results := make(chan artistData.Release)
	now := time.Now().UTC()
	for _, release := range releasesDetails {
		wg.Add(1)
		go t.prepareRelease(now, dbArtist, release, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	dbReleases := make([]artistData.Release, 0)
	for result := range results {
		dbReleases = append(dbReleases, result)
	}

	return dbReleases
}
