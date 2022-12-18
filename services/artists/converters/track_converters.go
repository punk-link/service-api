package converters

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"sort"
)

func ToTracksFromSpotify(tracks []releaseSpotifyPlatformModels.Track, artists map[string]artistData.Artist) []artistModels.Track {
	sort.SliceStable(tracks, func(i, j int) bool {
		return tracks[i].DiscNumber < tracks[j].DiscNumber && tracks[i].TrackNumber < tracks[j].TrackNumber
	})

	results := make([]artistModels.Track, len(tracks))
	for i, track := range tracks {
		results[i] = artistModels.Track{
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
