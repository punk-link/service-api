package repositories

import platformData "main/data/platforms"

type PlatformUrlRepository interface {
	CreatelsInBatches(err error, urls []platformData.PlatformReleaseUrl) error
	GetByReleaseId(err error, id int) ([]platformData.PlatformReleaseUrl, error)
	GetByReleaseIds(err error, ids []int) ([]platformData.PlatformReleaseUrl, error)
	UpdateInBatches(err error, urls []platformData.PlatformReleaseUrl) error
}
