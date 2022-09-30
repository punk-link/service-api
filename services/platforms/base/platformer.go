package base

import "main/models/platforms"

type Platformer interface {
	GetReleaseUrlsByUpc(upcContainers []platforms.UpcContainer) []platforms.UpcContainer
}
