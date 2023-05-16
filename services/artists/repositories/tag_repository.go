package repositories

import artistData "main/data/artists"

type TagRepository interface {
	Create(err error, tags *[]artistData.Tag) error
	Get(err error, normalizedNames []string) []artistData.Tag
	GetByReleaseId(err error, releaseId int) ([]artistData.Tag, error)
}
