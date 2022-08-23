package artists

import (
	"encoding/json"
	"main/data"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	commonModels "main/models/common"
	"main/models/labels"
	"main/models/spotify/releases"
	"main/services/common"
	"main/services/spotify"
	"sort"
	"sync"
	"time"

	"github.com/jackc/pgtype"
)

type ReleaseService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
}

func BuildReleaseService(logger *common.Logger, spotifyService *spotify.SpotifyService) *ReleaseService {
	return &ReleaseService{
		logger:         logger,
		spotifyService: spotifyService,
	}
}

func (t *ReleaseService) AddReleases(currentManager labels.ManagerContext, artists map[string]artistData.Artist, releases []releases.Release, timeStamp time.Time) ([]artistModels.Release, error) {
	var wg sync.WaitGroup
	chanResults := make(chan artistData.Release)
	for _, release := range releases {
		wg.Add(1)
		go t.buildReleaseFromSpotifyDetails(&wg, chanResults, timeStamp, artists, release)
	}

	go func() {
		wg.Wait()
		close(chanResults)
	}()

	dbReleases := make([]artistData.Release, 0)
	for result := range chanResults {
		dbReleases = append(dbReleases, result)
	}

	err := data.DB.CreateInBatches(&dbReleases, 50).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
		return make([]artistModels.Release, 0), err
	}

	return t.GetReleases(currentManager, dbReleases[0].Artists[0].Id)
}

func (t *ReleaseService) GetRelease(currentManager labels.ManagerContext, id int) artistModels.Release {
	return artistModels.Release{}
}

func (t *ReleaseService) GetReleases(currentManager labels.ManagerContext, artistId int) ([]artistModels.Release, error) {
	var dbReleases []artistData.Release
	err := data.DB.Model(&artistData.Release{}).Preload("Artists").Where("primary_artist_id = ?", artistId).Find(&dbReleases).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
		return make([]artistModels.Release, 0), err
	}

	return toReleases(dbReleases), nil
}

func toReleases(dbReleases []artistData.Release) []artistModels.Release {
	results := make([]artistModels.Release, len(dbReleases))
	for i, dbRelease := range dbReleases {

		results[i] = artistModels.Release{
			Id:           dbRelease.Id,
			Artists:      []artistModels.Artist{},
			ImageDetails: commonModels.ImageDetails{},
			Lable:        dbRelease.Label,
			Name:         dbRelease.Name,
			ReleaseDate:  dbRelease.ReleaseDate,
			TrackNumber:  dbRelease.TrackNumber,
			Tracks:       []artistModels.Track{},
			Type:         dbRelease.Type,
		}
	}

	return results
}

func (t *ReleaseService) buildReleaseFromSpotifyDetails(wg *sync.WaitGroup, results chan<- artistData.Release, timeStamp time.Time, artists map[string]artistData.Artist, release releases.Release) {
	defer wg.Done()

	imageDetails := t.getImageDetails(release)
	imageDetailsJson, err := json.Marshal(imageDetails)
	if err != nil {
		t.logger.LogError(err, "Can't serialize image details: '%s'", err.Error())
		return
	}

	imageDetailsJsonb := pgtype.JSONB{}
	err = imageDetailsJsonb.Set(imageDetailsJson)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	tracks := t.getTracks(release, artists)
	tracksJson, err := json.Marshal(tracks)
	if err != nil {
		t.logger.LogError(err, "Can't serialize track: '%s'", err.Error())
		return
	}

	releaseArtists := t.getReleaseArtists(release, artists)

	dbRelease := artistData.Release{
		Artists:         releaseArtists,
		Created:         timeStamp,
		ImageDetails:    string(imageDetailsJson),
		Label:           t.getLabelName(release),
		Name:            release.Name,
		PrimaryArtistId: releaseArtists[0].Id,
		ReleaseDate:     t.getReleaseDate(release),
		SpotifyId:       release.Id,
		TrackNumber:     release.TrackNumber,
		Tracks:          string(tracksJson),
		Type:            release.Type,
		Updated:         timeStamp,
	}

	results <- dbRelease
}

func (t *ReleaseService) GetMissingSpotifyReleases(dbArtist artistData.Artist) []releases.Release {
	dbReleaseIds := make(map[string]int, len(dbArtist.Releases))
	for _, release := range dbArtist.Releases {
		dbReleaseIds[release.SpotifyId] = 0
	}

	spotifyReleases := t.spotifyService.GetArtistReleases(dbArtist.SpotifyId)
	missingReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if _, isContains := dbReleaseIds[spotifyRelease.Id]; !isContains {
			missingReleaseSpotifyIds = append(missingReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	return t.spotifyService.GetReleasesDetails(missingReleaseSpotifyIds)
}

func (t *ReleaseService) getImageDetails(release releases.Release) commonModels.ImageDetails {
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

func (t *ReleaseService) getLabelName(release releases.Release) string {
	if release.Label == "" {
		return release.Artists[0].Name
	}

	return release.Label
}

func (t *ReleaseService) getReleaseArtists(release releases.Release, artists map[string]artistData.Artist) []artistData.Artist {
	featuredArtists := make([]artistData.Artist, 0)
	for _, track := range release.Tracks.Items {
		if len(track.Artists) < 2 {
			continue
		}

		trackArtists := make([]artistData.Artist, len(track.Artists)-1)
		for i := 1; i < len(track.Artists); i++ {
			trackArtists[i-1] = artists[track.Artists[i].Id]
		}

		featuredArtists = append(featuredArtists, trackArtists...)
	}

	results := make([]artistData.Artist, 1+len(featuredArtists))
	results[0] = artists[release.Artists[0].Id]
	results = append(results, featuredArtists...)

	return results
}

func (t *ReleaseService) getReleaseDate(release releases.Release) time.Time {
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

func (t *ReleaseService) getTracks(release releases.Release, artists map[string]artistData.Artist) []artistModels.Track {
	sort.SliceStable(release.Tracks.Items, func(i, j int) bool {
		return release.Tracks.Items[i].DiscNumber < release.Tracks.Items[j].DiscNumber && release.Tracks.Items[i].TrackNumber < release.Tracks.Items[j].TrackNumber
	})

	tracks := make([]artistModels.Track, len(release.Tracks.Items))
	for i, track := range release.Tracks.Items {
		trackArtists := make([]artistModels.Artist, len(track.Artists))
		for i, artist := range track.Artists {
			sort.SliceStable(artist.ImageDetails, func(i, j int) bool {
				return artist.ImageDetails[i].Height > artist.ImageDetails[j].Height
			})

			trackArtists[i] = artistModels.Artist{
				Id:           artists[artist.Id].Id,
				ImageDetails: commonModels.ImageDetails{},
				Name:         artist.Name,
			}
		}

		tracks[i] = artistModels.Track{
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
