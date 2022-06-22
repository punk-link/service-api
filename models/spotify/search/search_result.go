package search

type ArtistSearchResult struct {
	Artists Artists `json:"artists"`
}

type Artists struct {
	Items []Artist `json:"items"`
}

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
