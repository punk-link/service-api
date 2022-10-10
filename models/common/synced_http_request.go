package common

import "net/http"

type SyncedHttpRequest struct {
	HttpRequest *http.Request
	Sync        string
}
