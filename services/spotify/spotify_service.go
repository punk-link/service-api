package spotify

import (
	"errors"
	"main/models/spotify/search"
	"net/url"
)

func SearchArtist(query string) (search.ArtistSearchResult, error) {
	var response search.ArtistSearchResult

	if len(query) < 3 {
		return response, errors.New("the query contains less than 3 characters")
	}

	if err := MakeRequest("GET", "search?type=artist&limit=10&q="+url.QueryEscape(query), &response); err != nil {
		return response, err
	}

	return response, nil
}
