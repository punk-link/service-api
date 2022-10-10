package accessToken

import "time"

type SpotifyAccessTokenContainer struct {
	Token   string
	Expired time.Time
}
