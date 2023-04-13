package converters

import artistData "main/data/artists"

func ToDbReleaseTagRelation(releaseId int, tagId int) artistData.ReleaseTagRelation {
	return artistData.ReleaseTagRelation{
		ReleaseId: releaseId,
		TagId:     tagId,
	}
}
