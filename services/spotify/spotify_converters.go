package spotify

import (
	"main/models/spotify"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/responses"
)

func toArtistSearchResults(spotifyArtists []search.Artist) []responses.ArtistSearchResult {
	artists := make([]responses.ArtistSearchResult, len(spotifyArtists))
	for i, artist := range spotifyArtists {
		artists[i] = responses.ArtistSearchResult{
			ImageMetadata: toImageMetadataResponse(artist.ImageMetadata),
			Name:          artist.Name,
			SpotifyId:     artist.Id,
		}
	}

	return artists
}

func toImageMetadataResponse(metadatas []spotify.ImageMetadata) []responses.ImageMetadata {
	results := make([]responses.ImageMetadata, len(metadatas))
	for i, metadata := range metadatas {
		results[i] = responses.ImageMetadata{
			Height: metadata.Height,
			Url:    metadata.Url,
		}
	}

	return results
}

func toRelease(spotifyRelease releases.ArtistRelease) responses.ArtistRelease {
	return responses.ArtistRelease{
		SpotifyId:     spotifyRelease.Id,
		Artists:       toArtistSearchResults(spotifyRelease.Artists),
		ImageMetadata: toImageMetadataResponse(spotifyRelease.ImageMetadata),
		Lable:         spotifyRelease.Label,
		Name:          spotifyRelease.Name,
		ReleaseDate:   spotifyRelease.ReleaseDate,
		TrackNumber:   spotifyRelease.TrackNumber,
		Tracks:        toTracks(spotifyRelease.Tracks.Items),
		Type:          spotifyRelease.Type,
	}
}

func toReleases(spotifyReleases []releases.ArtistRelease) []responses.ArtistRelease {
	releases := make([]responses.ArtistRelease, len(spotifyReleases))
	for i, release := range spotifyReleases {
		releases[i] = toRelease(release)
	}

	return releases
}

func toTracks(spotifyTracks []releases.Track) []responses.Track {
	tracks := make([]responses.Track, len(spotifyTracks))
	for i, track := range spotifyTracks {
		tracks[i] = responses.Track{
			SpotifyId:       track.Id,
			Artists:         toArtistSearchResults(track.Artists),
			DiscNumber:      track.DiscNumber,
			DurationSeconds: track.DurationMilliseconds / 1000,
			IsExplicit:      track.IsExplicit,
			Name:            track.Name,
			TrackNumber:     track.TrackNumber,
		}
	}

	return tracks
}
