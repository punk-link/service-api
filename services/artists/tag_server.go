package artists

import artistModels "main/models/artists"

type TagServer interface {
	GetOrAdd(tagNames []string) []artistModels.Tag
}
