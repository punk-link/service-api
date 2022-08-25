package artists

import (
	"main/data"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	"main/models/labels"
	"main/models/spotify/releases"
	"main/services/artists/converters"
	"main/services/common"
	commonConverters "main/services/common/converters"
	"main/services/spotify"
	"sync"
	"time"
)

type ReleaseService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
}

func ConstructReleaseService(logger *common.Logger, spotifyService *spotify.SpotifyService) *ReleaseService {
	return &ReleaseService{
		logger:         logger,
		spotifyService: spotifyService,
	}
}

func (t *ReleaseService) Add(currentManager labels.ManagerContext, artists map[string]artistData.Artist, releases []releases.Release, timeStamp time.Time) error {
	dbReleases := t.buildDbReleases(artists, releases, timeStamp)

	err := data.DB.CreateInBatches(&dbReleases, 50).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}

func (t *ReleaseService) Get(currentManager labels.ManagerContext, artistId int) ([]artistModels.Release, error) {
	var dbReleases []artistData.Release
	err := data.DB.Model(&artistData.Release{}).Preload("Artists").Where("primary_artist_id = ?", artistId).Find(&dbReleases).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
		return make([]artistModels.Release, 0), err
	}

	releases, err := converters.ToReleases(dbReleases)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return releases, err
}

func (t *ReleaseService) GetMissingReleases(dbArtist artistData.Artist) []releases.Release {
	// TODO: replace dbArtist with ID

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

func (t *ReleaseService) GetOne(currentManager labels.ManagerContext, id int) artistModels.Release {
	return artistModels.Release{}
}

func (t *ReleaseService) buildDbReleases(artists map[string]artistData.Artist, releases []releases.Release, timeStamp time.Time) []artistData.Release {
	var wg sync.WaitGroup
	chanResults := make(chan artistData.Release)
	for _, release := range releases {
		wg.Add(1)
		go t.buildFromSpotify(&wg, chanResults, artists, release, timeStamp)
	}

	go func() {
		wg.Wait()
		close(chanResults)
	}()

	dbReleases := make([]artistData.Release, 0)
	for result := range chanResults {
		dbReleases = append(dbReleases, result)
	}

	return dbReleases
}

func (t *ReleaseService) buildFromSpotify(wg *sync.WaitGroup, results chan<- artistData.Release, artists map[string]artistData.Artist, release releases.Release, timeStamp time.Time) {
	defer wg.Done()

	imageDetails := commonConverters.ToImageDetailsFromSpotify(release.ImageDetails, release.Name)
	releaseArtists := t.getReleaseArtists(release, artists)
	tracks := converters.ToTracksFromSpotify(release.Tracks.Items, artists)

	dbArtist, err := converters.ToDbReleaseFromSpotify(release, releaseArtists, imageDetails, tracks, timeStamp)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return
	}

	results <- dbArtist
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
