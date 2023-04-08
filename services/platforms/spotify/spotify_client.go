package spotify

import "net/http"

type SpotifyClient interface {
	Request(params [][]string, method string, format string) []*http.Request
	RequestOne(method string, url string) (*http.Request, error)
}
