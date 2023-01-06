package converters

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"
	commonConverters "main/services/common/converters"

	presentationContracts "github.com/punk-link/presentation-contracts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToReleaseMessage(err error, release artistData.Release, artists []artistData.SlimArtist) (*presentationContracts.Release, error) {
	imageDetails, err := commonConverters.ToMessageFromJson(err, release.ImageDetails)
	slimArtists, err := ToSlimArtistMessages(err, artists)
	tracks, err := toTrackMessages(err, release.Tracks)
	if err != nil {
		return &presentationContracts.Release{}, err
	}

	return &presentationContracts.Release{
		Id:             int32(release.Id),
		ImageDetails:   imageDetails,
		Label:          release.Label,
		Name:           release.Name,
		ReleaseArtists: slimArtists,
		ReleaseDate:    timestamppb.New(release.ReleaseDate),
		Tracks:         tracks,
		Type:           release.Type,
	}, nil
}

func ToSlimReleaseMessages(err error, dbReleases []artistData.SlimRelease) ([]*presentationContracts.SlimRelease, error) {
	if err != nil {
		return make([]*presentationContracts.SlimRelease, 0), err
	}

	results := make([]*presentationContracts.SlimRelease, len(dbReleases))
	for i, dbRelease := range dbReleases {
		results[i] = toSlimReleaseMessage(dbRelease)
	}

	return results, nil
}

func toSlimReleaseMessage(dbRelease artistData.SlimRelease) *presentationContracts.SlimRelease {
	imageDetails, _ := commonConverters.ToMessageFromJson(nil, dbRelease.ImageDetails)

	return &presentationContracts.SlimRelease{
		Id:           int32(dbRelease.Id),
		ImageDetails: imageDetails,
		Name:         dbRelease.Name,
		ReleaseDate:  timestamppb.New(dbRelease.ReleaseDate),
		Type:         dbRelease.Type,
	}
}

func toTrackMessage(track artistModels.Track) *presentationContracts.Track {
	artists, _ := ToSlimArtistMessagesFromModels(nil, track.Artists)

	return &presentationContracts.Track{
		Artists:     artists,
		DiscNumber:  int32(track.DiscNumber),
		IsExplicit:  track.IsExplicit,
		Name:        track.Name,
		TrackNumber: int32(track.TrackNumber),
	}
}

func toTrackMessages(err error, json string) ([]*presentationContracts.Track, error) {
	if err != nil {
		return make([]*presentationContracts.Track, 0), err
	}

	tracks, err := getTracks(json)
	if err != nil {
		return make([]*presentationContracts.Track, 0), err
	}

	results := make([]*presentationContracts.Track, len(tracks))
	for i, track := range tracks {
		results[i] = toTrackMessage(track)
	}

	return results, nil
}
