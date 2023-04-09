package converters

import (
	commonConverters "main/converters/common"
	artistData "main/data/artists"
	platformData "main/data/platforms"
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToReleaseMessage(err error, release artistData.Release, artists []artistData.SlimArtist, platformUrls []platformData.PlatformReleaseUrl) (*presentationContracts.Release, error) {
	imageDetails, err := commonConverters.ToMessageFromJson(err, release.ImageDetails)
	slimArtists, err := ToSlimArtistMessages(err, artists)
	tracks, err := toTrackMessages(err, release.Tracks)
	platformUrlMessages, err := toPlatformUrlMessages(err, platformUrls)
	if err != nil {
		return &presentationContracts.Release{}, err
	}

	return &presentationContracts.Release{
		Id: int32(release.Id),
		ArtistStats: &presentationContracts.ArtistStats{
			CompilationNumber: int32(3),
			SoleReleaseNumber: int32(1),
		},
		Description:    "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
		ImageDetails:   imageDetails,
		Label:          release.Label,
		Name:           release.Name,
		PlatformUrls:   platformUrlMessages,
		ReleaseArtists: slimArtists,
		ReleaseDate:    timestamppb.New(release.ReleaseDate),
		Tags:           []string{"indie", "post-punk"},
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
