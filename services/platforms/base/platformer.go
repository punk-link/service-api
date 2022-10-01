package base

import "main/models/platforms"

type Platformer interface {
	GetPlatformName() string
	GetReleaseUrlsByUpc(upcContainers []platforms.UpcContainer) []platforms.UrlResultContainer
}
