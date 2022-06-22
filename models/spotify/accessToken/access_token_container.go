package accessToken

import "time"

type SpotifyAccessTokenContainer struct {
	Token     string
	ExpiresAt time.Time
}
