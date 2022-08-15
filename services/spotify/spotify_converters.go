package spotify

import (
	"main/models/artists"
	"main/models/spotify"
	"main/models/spotify/releases"
	"main/models/spotify/search"
)

func ToArtistSearchResults(spotifyArtists []search.Artist) []artists.ArtistSearchResult {
	results := make([]artists.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		results[i] = artists.ArtistSearchResult{
			ImageMetadata: toImageMetadataResponse(artist.ImageMetadata),
			Name:          artist.Name,
			SpotifyId:     artist.Id,
		}
	}

	return results
}

func toImageMetadataResponse(metadatas []spotify.ImageMetadata) []artists.ImageMetadata {
	results := make([]artists.ImageMetadata, len(metadatas))
	for i, metadata := range metadatas {
		results[i] = artists.ImageMetadata{
			Height: metadata.Height,
			Url:    metadata.Url,
		}
	}

	return results
}

func toRelease(spotifyRelease releases.Release) artists.Release {
	return artists.Release{
		SpotifyId:     spotifyRelease.Id,
		Artists:       ToArtistSearchResults(spotifyRelease.Artists),
		ImageMetadata: toImageMetadataResponse(spotifyRelease.ImageMetadata),
		Lable:         spotifyRelease.Label,
		Name:          spotifyRelease.Name,
		ReleaseDate:   spotifyRelease.ReleaseDate,
		TrackNumber:   spotifyRelease.TrackNumber,
		Tracks:        toTracks(spotifyRelease.Tracks.Items),
		Type:          spotifyRelease.Type,
	}
}

func toReleases(spotifyReleases []releases.Release) []artists.Release {
	releases := make([]artists.Release, len(spotifyReleases))
	for i, release := range spotifyReleases {
		releases[i] = toRelease(release)
	}

	return releases
}

func toTracks(spotifyTracks []releases.Track) []artists.Track {
	tracks := make([]artists.Track, len(spotifyTracks))
	for i, track := range spotifyTracks {
		tracks[i] = artists.Track{
			SpotifyId:       track.Id,
			Artists:         ToArtistSearchResults(track.Artists),
			DiscNumber:      track.DiscNumber,
			DurationSeconds: track.DurationMilliseconds / 1000,
			IsExplicit:      track.IsExplicit,
			Name:            track.Name,
			TrackNumber:     track.TrackNumber,
		}
	}

	return tracks
}
