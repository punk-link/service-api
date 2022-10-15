package platforms

import (
	"main/data"
	platformData "main/data/platforms"
	"main/helpers"

	"github.com/punk-link/logger"
)

func createDbPlatformReleaseUrlsInBatches(logger *logger.Logger, err error, urls []platformData.PlatformReleaseUrl) error {
	if err != nil {
		return err
	}

	err = data.DB.CreateInBatches(&urls, CREATE_PLATFORM_RELEASE_URLS_BATCH_SIZE).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbPlatformReleaseUrlsByReleaseId(logger *logger.Logger, err error, id int) ([]platformData.PlatformReleaseUrl, error) {
	return getDbPlatformReleaseUrlsByReleaseIds(logger, err, []int{id})
}

func getDbPlatformReleaseUrlsByReleaseIds(logger *logger.Logger, err error, ids []int) ([]platformData.PlatformReleaseUrl, error) {
	if err != nil {
		return make([]platformData.PlatformReleaseUrl, 0), err
	}

	var results []platformData.PlatformReleaseUrl
	err = data.DB.
		Where("release_id in (?)", ids).
		Find(&results).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return results, err
}

func updateDbPlatformReleaseUrlsInBatches(logger *logger.Logger, err error, urls []platformData.PlatformReleaseUrl) error {
	if err != nil || len(urls) == 0 {
		return err
	}

	for _, url := range urls {
		innerErr := data.DB.Model(&platformData.PlatformReleaseUrl{}).
			Where("id = ?", url.Id).
			Updates(url).
			Error

		err = helpers.CombineErrors(err, innerErr)
	}

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_PLATFORM_RELEASE_URLS_BATCH_SIZE int = 200
