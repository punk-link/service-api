package repositories

import artistData "main/data/artists"

type SocialNetworkRepository interface {
	Get(err error, artistId int) []artistData.ArtistSocialNetwork
}
