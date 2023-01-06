package helpers

import "fmt"

func BuildCacheKey(bucket string, key any) string {
	return fmt.Sprintf("%s::%s", bucket, key)
}
