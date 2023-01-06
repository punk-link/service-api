package converters

import (
	artistModels "main/models/artists"
	commonModels "main/models/common"

	presentationContracts "github.com/punk-link/presentation-contracts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToArtistMessage(artist artistModels.Artist) *presentationContracts.Artist {
	return &presentationContracts.Artist{
		Id:           int32(artist.Id),
		LabelId:      int32(artist.LabelId),
		Name:         artist.Name,
		ImageDetails: toImageDetailsMessage(artist.ImageDetails),
		Releases:     toReleaseMessages(artist.Releases),
	}
}

func toImageDetailsMessage(imageDetails commonModels.ImageDetails) *presentationContracts.ImageDetails {
	return &presentationContracts.ImageDetails{
		AltText: imageDetails.AltText,
		Height:  int32(imageDetails.Width),
		Url:     imageDetails.Url,
		Width:   int32(imageDetails.Width),
	}
}

func ToReleaseMessage(release artistModels.Release) *presentationContracts.Release {
	return &presentationContracts.Release{
		Id:               int32(release.Id),
		FeaturingArtists: toSlimArtistMessages(release.FeaturingArtists),
		ImageDetails:     toImageDetailsMessage(release.ImageDetails),
		Label:            release.Label,
		Name:             release.Name,
		ReleaseArtists:   toSlimArtistMessages(release.ReleaseArtists),
		ReleaseDate:      timestamppb.New(release.ReleaseDate),
		TrackNumber:      int32(release.TrackNumber),
		Tracks:           toTrackMessages(release.Tracks),
		Type:             release.Type,
	}
}

func toReleaseMessages(releases []artistModels.Release) []*presentationContracts.Release {
	results := make([]*presentationContracts.Release, len(releases))
	for i, release := range releases {
		results[i] = ToReleaseMessage(release)
	}

	return results
}

func toSlimArtistMessage(artist artistModels.Artist) *presentationContracts.SlimArtist {
	return &presentationContracts.SlimArtist{
		Id:           int32(artist.Id),
		LabelId:      int32(artist.LabelId),
		Name:         artist.Name,
		ImageDetails: toImageDetailsMessage(artist.ImageDetails),
	}
}

func toSlimArtistMessages(artists []artistModels.Artist) []*presentationContracts.SlimArtist {
	results := make([]*presentationContracts.SlimArtist, len(artists))
	for i, artist := range artists {
		results[i] = toSlimArtistMessage(artist)
	}

	return results
}

func toTrackMessage(track artistModels.Track) *presentationContracts.Track {
	return &presentationContracts.Track{
		Artists:     toSlimArtistMessages(track.Artists),
		DiscNumber:  int32(track.DiscNumber),
		IsExplicit:  track.IsExplicit,
		Name:        track.Name,
		TrackNumber: int32(track.TrackNumber),
	}
}

func toTrackMessages(tracks []artistModels.Track) []*presentationContracts.Track {
	results := make([]*presentationContracts.Track, len(tracks))
	for i, track := range tracks {
		results[i] = toTrackMessage(track)
	}

	return results
}
