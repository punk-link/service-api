package spotify

import (
	"errors"
	"main/models/spotify"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/responses"
	"net/url"
	"strconv"
)

func GetArtistRelease(spotifyId string) (responses.ArtistRelease, error) {
	var result responses.ArtistRelease

	urlPattern := "albums/" + spotifyId
	var spotifyRelease releases.ArtistRelease
	if err := MakeRequest("GET", urlPattern, &spotifyRelease); err != nil {
		return result, err
	}

	return toRelease(spotifyRelease), nil
}

func GetArtistReleases(spotifyId string) ([]responses.ArtistRelease, error) {
	spotifyReleases := getReleases(spotifyId)
	return toReleases(spotifyReleases), nil
}

func SearchArtist(query string) ([]responses.ArtistSearchResult, error) {
	const minimalQueryLength int = 3
	const urlPattern string = "search?type=artist&limit=10&q="
	var result []responses.ArtistSearchResult

	if len(query) < minimalQueryLength {
		return result, errors.New("the query contains less than 3 characters")
	}

	var spotifyResponse search.ArtistSearchResult
	if err := MakeRequest("GET", urlPattern+url.QueryEscape(query), &spotifyResponse); err != nil {
		return result, err
	}

	result = toArtistSearchResults(spotifyResponse.Artists.Items)
	return result, nil
}

func getReleases(spotifyId string) []releases.ArtistRelease {
	const queryLimit int = 20

	var spotifyResponse releases.ArtistReleaseResult
	offset := 0
	for {
		urlPattern := "artists/" + spotifyId + "/albums?limit=" + strconv.Itoa(queryLimit) + "&offset=" + strconv.Itoa(offset)
		offset = offset + queryLimit

		var tmpResponse releases.ArtistReleaseResult
		if err := MakeRequest("GET", urlPattern, &tmpResponse); err != nil {
			// TODO: log an error
			continue
		}

		spotifyResponse.Items = append(spotifyResponse.Items, tmpResponse.Items...)
		if tmpResponse.Next == "" {
			break
		}
	}

	return spotifyResponse.Items
}

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

func toReleases(spotifyReleases []releases.ArtistRelease) []responses.ArtistRelease {
	releases := make([]responses.ArtistRelease, len(spotifyReleases))
	for i, release := range spotifyReleases {
		releases[i] = toRelease(release)
	}

	return releases
}
