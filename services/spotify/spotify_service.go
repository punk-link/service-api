package spotify

import (
	"fmt"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/responses"
	spotifyConverters "main/utils/converters/spotify"
	"net/url"
)

func GetArtistRelease(spotifyId string) responses.ArtistRelease {
	var result responses.ArtistRelease

	var spotifyRelease releases.ArtistRelease
	if err := MakeRequest("GET", fmt.Sprintf("albums/%s", spotifyId), &spotifyRelease); err != nil {
		fmt.Println(err)
		return result
	}

	return spotifyConverters.ToRelease(spotifyRelease)
}

func GetArtistReleases(spotifyId string) []responses.ArtistRelease {
	spotifyReleases := getReleases(spotifyId)
	return spotifyConverters.ToReleases(spotifyReleases)
}

func SearchArtist(query string) []responses.ArtistSearchResult {
	var result []responses.ArtistSearchResult

	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return result
	}

	var spotifyArtistSearchResults search.ArtistSearchResult
	if err := MakeRequest("GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults); err != nil {
		fmt.Println(err)
		return result
	}

	return spotifyConverters.ToArtistSearchResults(spotifyArtistSearchResults.Artists.Items)
}

func getReleases(spotifyId string) []releases.ArtistRelease {
	const queryLimit int = 20

	var spotifyResponse releases.ArtistReleaseResult
	offset := 0
	for {
		urlPattern := fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, queryLimit, offset)
		offset = offset + queryLimit

		var tmpResponse releases.ArtistReleaseResult
		if err := MakeRequest("GET", urlPattern, &tmpResponse); err != nil {
			fmt.Println(err)
			continue
		}

		spotifyResponse.Items = append(spotifyResponse.Items, tmpResponse.Items...)
		if tmpResponse.Next == "" {
			break
		}
	}

	return spotifyResponse.Items
}
