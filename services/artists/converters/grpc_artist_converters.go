package converters

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
)

func ToArtistMessage(err error, artist artistData.SlimArtist, releaseStats artistModels.ArtistReleaseStats, socialNetworks []artistModels.SocialNetwork, presentationConfig artistModels.PresentationConfig) (*presentationContracts.Artist, error) {
	presentationConfigMessage, err := ToPresentationConfigMessage(err, presentationConfig)
	releaseStatsMessage, err := ToReleaseStatsMessage(err, releaseStats)
	if err != nil {
		return &presentationContracts.Artist{}, err
	}

	return &presentationContracts.Artist{
		Id:                 int32(artist.Id),
		Name:               artist.Name,
		PresentationConfig: presentationConfigMessage,
		Releases:           nil,
		ReleaseStats:       releaseStatsMessage,
		SocialNetworks:     ToSocialNetworksMap(socialNetworks),
	}, nil
}

func ToSlimArtistMessages(err error, artists []artistData.SlimArtist) ([]*presentationContracts.SlimArtist, error) {
	if err != nil {
		return make([]*presentationContracts.SlimArtist, 0), err
	}

	results := make([]*presentationContracts.SlimArtist, len(artists))
	for i, artist := range artists {
		results[i] = toSlimArtistMessage(artist)
	}

	return results, nil
}

func ToSlimArtistMessagesFromModels(err error, artists []artistModels.Artist) ([]*presentationContracts.SlimArtist, error) {
	if err != nil {
		return make([]*presentationContracts.SlimArtist, 0), err
	}

	dbArtists := make([]artistData.SlimArtist, len(artists))
	for i, artist := range artists {
		dbArtists[i] = artistData.SlimArtist{
			Id:   artist.Id,
			Name: artist.Name,
		}
	}

	return ToSlimArtistMessages(nil, dbArtists)
}

func toSlimArtistMessage(artist artistData.SlimArtist) *presentationContracts.SlimArtist {
	return &presentationContracts.SlimArtist{
		Id:   int32(artist.Id),
		Name: artist.Name,
	}
}
