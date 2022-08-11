package consul

import "time"

type LocalCacheContainer struct {
	Expired time.Time
	Value   interface{}
}
