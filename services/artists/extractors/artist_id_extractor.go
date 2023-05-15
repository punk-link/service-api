package extractors

import artistData "main/data/artists"

type ArtistIdExtractor interface {
	Extract(releases []artistData.Release) []int
	ExtractFromOne(release artistData.Release) []int
}
