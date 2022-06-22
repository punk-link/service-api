package spotify

import (
	"errors"
	"main/models/spotify/search"
	"main/responses"
	"net/url"
)

func SearchArtist(query string) (responses.ArtistSearchResponse, error) {
	const minimalQueryLength int = 3
	const urlPattern string = "search?type=artist&limit=10&q="
	var result responses.ArtistSearchResponse

	if len(query) < minimalQueryLength {
		return result, errors.New("the query contains less than 3 characters")
	}

	var response search.ArtistSearchResult
	if err := MakeRequest("GET", urlPattern+url.QueryEscape(query), &response); err != nil {
		return result, err
	}

	artistNumber := len(response.Artists.Items)
	artists := make([]responses.ArtistSearchResult, artistNumber)
	for i := 0; i < artistNumber; i++ {
		artists[i] = responses.ArtistSearchResult{
			Name:      response.Artists.Items[i].Name,
			SpotifyId: response.Artists.Items[i].Id,
		}
	}

	result.Items = artists

	return result, nil
}
