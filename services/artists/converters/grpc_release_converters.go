package converters

import (
	commonConverters "main/converters/common"
	artistData "main/data/artists"
	platformData "main/data/platforms"
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToReleaseMessage(err error, release artistData.Release, artists []artistData.SlimArtist, platformUrls []platformData.PlatformReleaseUrl, tags []string, releaseStats artistModels.ArtistReleaseStats, presentationConfig artistModels.PresentationConfig) (*presentationContracts.Release, error) {
	imageDetails, err := commonConverters.ToMessageFromJson(err, release.ImageDetails)
	platformUrlMessages, err := toPlatformUrlMessages(err, platformUrls)
	presentationConfigMessage, err := ToPresentationConfigMessage(err, presentationConfig)
	releaseStatsMessage, err := ToReleaseStatsMessage(err, releaseStats)
	slimArtists, err := ToSlimArtistMessages(err, artists)
	tracks, err := toTrackMessages(err, release.Tracks)
	if err != nil {
		return &presentationContracts.Release{}, err
	}

	return &presentationContracts.Release{
		Id:                 int32(release.Id),
		Description:        release.Description,
		ImageDetails:       imageDetails,
		Label:              release.Label,
		Name:               release.Name,
		PlatformUrls:       platformUrlMessages,
		PresentationConfig: presentationConfigMessage,
		ReleaseArtists:     slimArtists,
		ReleaseDate:        timestamppb.New(release.ReleaseDate),
		ReleaseStats:       releaseStatsMessage,
		Tags:               tags,
		Tracks:             tracks,
		Type:               release.Type,
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

func toPlatformUrlMessage(platformUrl platformData.PlatformReleaseUrl) *presentationContracts.PlatformUrl {
	return &presentationContracts.PlatformUrl{
		PlatformId: platformUrl.PlatformName,
		Url:        platformUrl.Url,
	}
}

func toPlatformUrlMessages(err error, platformUrls []platformData.PlatformReleaseUrl) ([]*presentationContracts.PlatformUrl, error) {
	if err != nil {
		return make([]*presentationContracts.PlatformUrl, 0), err
	}

	results := make([]*presentationContracts.PlatformUrl, len(platformUrls))
	for i, url := range platformUrls {
		results[i] = toPlatformUrlMessage(url)
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
