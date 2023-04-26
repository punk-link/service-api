package artists

type Tag struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
}
