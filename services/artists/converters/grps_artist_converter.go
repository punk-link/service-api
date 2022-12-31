package converters

import (
	presentationGrpcs "main/grpc/presentations"
	artistModels "main/models/artists"
	commonModels "main/models/common"
)

func ToArtistResponseMessage(artist artistModels.Artist) *presentationGrpcs.ArtistResponse {
	return &presentationGrpcs.ArtistResponse{
		Id:           int32(artist.Id),
		LabelId:      int32(artist.LabelId),
		Name:         artist.Name,
		ImageDetails: toImageDetailsMessage(artist.ImageDetails),
		Releases:     toReleaseMessages(artist.Releases),
	}
}

func toImageDetailsMessage(imageDetails commonModels.ImageDetails) *presentationGrpcs.ImageDetails {
	return &presentationGrpcs.ImageDetails{
		AltText: imageDetails.AltText,
		Height:  int32(imageDetails.Width),
		Url:     imageDetails.Url,
		Width:   int32(imageDetails.Width),
	}
}

func toReleaseMessage(release artistModels.Release) *presentationGrpcs.Release {
	return &presentationGrpcs.Release{
		Id:               int32(release.Id),
		FeaturingArtists: toSlimArtistMessages(release.FeaturingArtists),
		ImageDetails:     toImageDetailsMessage(release.ImageDetails),
		Label:            release.Label,
		Name:             release.Name,
		ReleaseArtists:   toSlimArtistMessages(release.ReleaseArtists),
		TrackNumber:      int32(release.TrackNumber),
		Tracks:           toTrackMessages(release.Tracks),
		Type:             release.Type,
	}
}

func toReleaseMessages(releases []artistModels.Release) []*presentationGrpcs.Release {
	results := make([]*presentationGrpcs.Release, len(releases))
	for i, release := range releases {
		results[i] = toReleaseMessage(release)
	}

	return results
}

func toSlimArtistMessage(artist artistModels.Artist) *presentationGrpcs.SlimArtist {
	return &presentationGrpcs.SlimArtist{
		Id:           int32(artist.Id),
		LabelId:      int32(artist.LabelId),
		Name:         artist.Name,
		ImageDetails: toImageDetailsMessage(artist.ImageDetails),
	}
}

func toSlimArtistMessages(artists []artistModels.Artist) []*presentationGrpcs.SlimArtist {
	results := make([]*presentationGrpcs.SlimArtist, len(artists))
	for i, artist := range artists {
		results[i] = toSlimArtistMessage(artist)
	}

	return results
}

func toTrackMessage(track artistModels.Track) *presentationGrpcs.Track {
	return &presentationGrpcs.Track{
		Artists:     toSlimArtistMessages(track.Artists),
		DiscNumber:  int32(track.DiscNumber),
		IsExplicit:  track.IsExplicit,
		Name:        track.Name,
		TrackNumber: int32(track.TrackNumber),
	}
}

func toTrackMessages(tracks []artistModels.Track) []*presentationGrpcs.Track {
	results := make([]*presentationGrpcs.Track, len(tracks))
	for i, track := range tracks {
		results[i] = toTrackMessage(track)
	}

	return results
}
