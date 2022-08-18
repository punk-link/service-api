package spotify

import (
	"main/models/artists"
	"main/models/common"
	"main/models/spotify"
	"main/models/spotify/releases"
	"main/models/spotify/search"
)

func ToArtist(spotifyArtists []search.Artist) []artists.Artist {
	results := make([]artists.Artist, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = artists.Artist{
			ImageDetails: toImageDetailsaResponse(artist.ImageDetails)[0],
			Name:         artist.Name,
			Id:           0,
		}
	}

	return results
}

func ToArtistSearchResults(spotifyArtists []search.Artist) []artists.ArtistSearchResult {
	results := make([]artists.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		imageDetails := common.ImageDetails{}
		if 0 < len(artist.ImageDetails) {
			imageDetails = toImageDetailsaResponse(artist.ImageDetails)[0]
		}

		results[i] = artists.ArtistSearchResult{
			ImageDetails: imageDetails,
			Name:         artist.Name,
			SpotifyId:    artist.Id,
		}
	}

	return results
}

func toImageDetailsaResponse(imageDetails []spotify.ImageDetails) []common.ImageDetails {
	results := make([]common.ImageDetails, len(imageDetails))
	for i, details := range imageDetails {
		results[i] = common.ImageDetails{
			Height: details.Height,
			Url:    details.Url,
		}
	}

	return results
}

func toRelease(spotifyRelease releases.Release) artists.Release {
	return artists.Release{
		SpotifyId:    spotifyRelease.Id,
		Artists:      ToArtist(spotifyRelease.Artists),
		ImageDetails: toImageDetailsaResponse(spotifyRelease.ImageDetails),
		Lable:        spotifyRelease.Label,
		Name:         spotifyRelease.Name,
		ReleaseDate:  spotifyRelease.ReleaseDate,
		TrackNumber:  spotifyRelease.TrackNumber,
		Tracks:       toTracks(spotifyRelease.Tracks.Items),
		Type:         spotifyRelease.Type,
	}
}

func toTracks(spotifyTracks []releases.Track) []artists.Track {
	tracks := make([]artists.Track, len(spotifyTracks))
	for i, track := range spotifyTracks {
		tracks[i] = artists.Track{
			SpotifyId:       track.Id,
			Artists:         ToArtist(track.Artists),
			DiscNumber:      track.DiscNumber,
			DurationSeconds: track.DurationMilliseconds / 1000,
			IsExplicit:      track.IsExplicit,
			Name:            track.Name,
			TrackNumber:     track.TrackNumber,
		}
	}

	return tracks
}
