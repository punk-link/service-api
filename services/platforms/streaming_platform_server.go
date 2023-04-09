package platforms

import platformModels "main/models/platforms"

type StreamingPlatformServer interface {
	Get(releaseId int) ([]platformModels.PlatformReleaseUrl, error)
	ProcessPlatforeUrlResults()
	PublishPlatforeUrlRequests()
}
