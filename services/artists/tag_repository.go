package artists

import artistData "main/data/artists"

type TagRepository interface {
	Create(err error, tags *[]artistData.Tag) error
	Get(err error, normalizedNames []string) []artistData.Tag
}
