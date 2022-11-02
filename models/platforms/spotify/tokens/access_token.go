package tokens

type SpotifyAccessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}
