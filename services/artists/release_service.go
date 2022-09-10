package artists

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"
	"main/models/labels"
	"main/models/spotify/releases"
	"main/services/artists/converters"
	"main/services/common"
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
	orderReleasesChronologically(dbReleases)

	return createDbReleasesInBatches(t.logger, nil, &dbReleases)
}

func (t *ReleaseService) Get(artistId int) ([]artistModels.Release, error) {
	dbReleases, err := getDbReleasesByArtistId(t.logger, nil, artistId)
	artists, err := t.getReleasesArtists(err, dbReleases)
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	releases, err := converters.ToReleases(dbReleases, artists)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return releases, err
}

func (t *ReleaseService) GetMissingReleases(artistId int, artistSpotifyId string) ([]releases.Release, error) {
	dbReleases, err := getDbReleasesByArtistId(t.logger, nil, artistId)
	missingReleaseSpotifyIds, err := t.getMissingReleasesSpotifyIds(err, dbReleases, artistSpotifyId)
	if err != nil {
		return make([]releases.Release, 0), err
	}

	return t.spotifyService.GetReleasesDetails(missingReleaseSpotifyIds), nil
}

func (t *ReleaseService) GetOne(id int) artistModels.Release {
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

	dbArtist, err := converters.ToDbReleaseFromSpotify(release, artists, timeStamp)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return
	}

	results <- dbArtist
}

func (t *ReleaseService) getMissingReleasesSpotifyIds(err error, dbReleases []artistData.Release, artistSpotifyId string) ([]string, error) {
	if err != nil {
		return make([]string, 0), err
	}

	dbReleaseIds := make(map[string]int, len(dbReleases))
	for _, release := range dbReleases {
		dbReleaseIds[release.SpotifyId] = 0
	}

	spotifyReleases := t.spotifyService.GetArtistReleases(artistSpotifyId)
	missingReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if _, isContains := dbReleaseIds[spotifyRelease.Id]; !isContains {
			missingReleaseSpotifyIds = append(missingReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	return missingReleaseSpotifyIds, nil
}

func (t *ReleaseService) getReleasesArtists(err error, releases []artistData.Release) (map[int]artistModels.Artist, error) {
	if err != nil {
		return make(map[int]artistModels.Artist, 0), err
	}

	artistIds := getArtistsIdsFromDbReleases(t.logger, releases)
	artists, err := getDbArtists(t.logger, err, artistIds)

	results := make(map[int]artistModels.Artist, len(artists))
	for _, dbArtist := range artists {
		artist, err := converters.ToArtist(dbArtist, make([]artistModels.Release, 0))
		if err != nil {
			results[artist.Id] = artist
		}
	}

	return results, err
}
