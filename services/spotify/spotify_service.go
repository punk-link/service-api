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

func GetArtistReleases(spotifyId string) (responses.ArtistReleasesResponse, error) {
	result := responses.ArtistReleasesResponse{
		ArtistSpotifyId: spotifyId,
	}

	spotifyReleases := getReleases(spotifyId)
	result.Items = toReleases(spotifyReleases)

	return result, nil
}

func SearchArtist(query string) ([]responses.ArtistSearchResult, error) {
	const minimalQueryLength int = 3
	const urlPattern string = "search?type=artist&limit=10&q="
	var result []responses.ArtistSearchResult

	if len(query) < minimalQueryLength {
		return result, errors.New("the query contains less than 3 characters")
	}

	var response search.ArtistSearchResult
	if err := MakeRequest("GET", urlPattern+url.QueryEscape(query), &response); err != nil {
		return result, err
	}

	result = toArtistSearchResults(response.Artists.Items)
	return result, nil
}

func getReleases(spotifyId string) []releases.ArtistRelease {
	const queryLimit int = 20

	var response releases.ArtistReleaseResult
	offset := 0
	for {
		urlPattern := "artists/" + spotifyId + "/albums?limit=" + strconv.Itoa(queryLimit) + "&offset=" + strconv.Itoa(offset)
		offset = offset + queryLimit

		var tmpResponse releases.ArtistReleaseResult
		if err := MakeRequest("GET", urlPattern, &tmpResponse); err != nil {
			// TODO: log an error
			continue
		}

		response.Items = append(response.Items, tmpResponse.Items...)
		if tmpResponse.Next == "" {
			break
		}
	}

	return response.Items
}

func toArtistSearchResults(spotifyArtists []search.Artist) []responses.ArtistSearchResult {
	artistCount := len(spotifyArtists)
	artists := make([]responses.ArtistSearchResult, artistCount)
	for i := 0; i < artistCount; i++ {
		artists[i] = responses.ArtistSearchResult{
			ImageMetadata: toImageMetadataResponse(spotifyArtists[i].ImageMetadata),
			Name:          spotifyArtists[i].Name,
			SpotifyId:     spotifyArtists[i].Id,
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

func toReleases(spotifyReleases []releases.ArtistRelease) []responses.ArtistRelease {
	releaseNumber := len(spotifyReleases)
	releases := make([]responses.ArtistRelease, releaseNumber)
	for i := 0; i < releaseNumber; i++ {
		releases[i] = responses.ArtistRelease{
			SpotifyId:     spotifyReleases[i].Id,
			ImageMetadata: toImageMetadataResponse(spotifyReleases[i].ImageMetadata),
			Name:          spotifyReleases[i].Name,
			ReleaseDate:   spotifyReleases[i].ReleaseDate,
			TrackNumber:   spotifyReleases[i].TrackNumber,
			Type:          spotifyReleases[i].Type,
		}
	}

	return releases
}
