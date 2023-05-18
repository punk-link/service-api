package converters

import (
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
)

func ToSocialNetworksMap(socialNetworks []artistModels.SocialNetwork) []*presentationContracts.ArtistSocialNetwork {
	results := make([]*presentationContracts.ArtistSocialNetwork, len(socialNetworks))
	for i, network := range socialNetworks {
		results[i] = toSocialNetworkMap(network)
	}

	return results
}

func toSocialNetworkMap(socialNetwork artistModels.SocialNetwork) *presentationContracts.ArtistSocialNetwork {
	return &presentationContracts.ArtistSocialNetwork{
		Id:  socialNetwork.Id,
		Url: socialNetwork.Url,
	}
}
