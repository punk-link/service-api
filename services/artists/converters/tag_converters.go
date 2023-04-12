package converters

import artistData "main/data/artists"

func ToDbTag(name string, normalizedName string) artistData.Tag {
	return artistData.Tag{
		Name:           name,
		NormalizedName: normalizedName,
	}
}
