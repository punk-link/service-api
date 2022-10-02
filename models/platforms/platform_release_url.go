package platforms

type PlatformReleaseUrl struct {
	Id           int    `json:"id"`
	PlatformName string `json:"platformName"`
	ReleaseId    int    `json:"releaseId"`
	Url          string `json:"url"`
}
