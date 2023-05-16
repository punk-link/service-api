package artists

import artistModels "main/models/artists"

type TagServer interface {
	GetNames(releaseId int) []string
	GetOrAdd(tagNames []string) []artistModels.Tag
}
