package converters

import (
	data "main/data/artists"
	models "main/models/artists"
	spotifyReleases "main/models/platforms/spotify/releases"
	"sort"
)

func ToTracksFromSpotify(tracks []spotifyReleases.Track, artists map[string]data.Artist) []models.Track {
	sort.SliceStable(tracks, func(i, j int) bool {
		return tracks[i].DiscNumber < tracks[j].DiscNumber && tracks[i].TrackNumber < tracks[j].TrackNumber
	})

	results := make([]models.Track, len(tracks))
	for i, track := range tracks {
		results[i] = models.Track{
			Artists:         ToArtistsFromSpotifyTrack(track, artists),
			DiscNumber:      track.DiscNumber,
			DurationSeconds: track.DurationMilliseconds / 1000,
			IsExplicit:      track.IsExplicit,
			Name:            track.Name,
			SpotifyId:       track.Id,
			TrackNumber:     track.TrackNumber,
		}
	}

	return results
}
