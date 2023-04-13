package artists

import (
	dataStructures "main/data-structures"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	"main/services/artists/converters"
	"strings"
	"sync"
	"time"

	cacheManager "github.com/punk-link/cache-manager"
	"github.com/samber/do"
)

type TagService struct {
	mutex         sync.Mutex
	tagCache      cacheManager.CacheManager[map[string]artistModels.Tag]
	tagRepository TagRepository
}

func NewTagService(injector *do.Injector) (TagServer, error) {
	tagCache := do.MustInvoke[cacheManager.CacheManager[map[string]artistModels.Tag]](injector)
	tagRepository := do.MustInvoke[TagRepository](injector)

	return &TagService{
		tagCache:      tagCache,
		tagRepository: tagRepository,
	}, nil
}

func (t *TagService) GetOrAdd(tagNames []string) []artistModels.Tag {
	if len(tagNames) == 0 {
		return make([]artistModels.Tag, 0)
	}

	cache, isCached := t.tagCache.TryGet(TAG_CACHE_SLUG)
	if !isCached {
		cache = make(map[string]artistModels.Tag)
	}

	tags := make([]artistModels.Tag, len(tagNames))
	dbTags := make([]artistData.Tag, 0)
	normalizedNames := make([]string, 0)
	for _, name := range tagNames {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", ""))

		tag, isExist := cache[normalizedName]
		if isExist {
			tags = append(tags, tag)
			continue
		}

		dbTags = append(dbTags, converters.ToDbTag(name, normalizedName))

		normalizedNames = append(normalizedNames, normalizedName)
	}

	if len(dbTags) == 0 {
		return tags
	}

	existingDbTags := t.tagRepository.Get(nil, normalizedNames)
	existingTagNames := make([]string, len(existingDbTags))
	for i, tag := range existingDbTags {
		existingTagNames[i] = tag.NormalizedName

		existingTag := artistModels.Tag{
			Id:             tag.Id,
			Name:           tag.Name,
			NormalizedName: tag.NormalizedName,
		}
		tags = append(tags, existingTag)

		cache[existingTag.NormalizedName] = existingTag
	}
	existingTagNameSet := dataStructures.MakeHashSet(existingTagNames)

	newDbTags := make([]artistData.Tag, 0)
	for _, tag := range dbTags {
		if existingTagNameSet.Contains(tag.NormalizedName) {
			newDbTags = append(newDbTags, tag)
		}
	}

	err := t.tagRepository.Create(nil, &newDbTags)
	if err != nil {
		return tags
	}

	for _, tag := range newDbTags {
		newTag := artistModels.Tag{
			Id:             tag.Id,
			Name:           tag.Name,
			NormalizedName: tag.NormalizedName,
		}
		tags = append(tags, newTag)

		cache[newTag.NormalizedName] = newTag
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	actualCache, _ := t.tagCache.TryGet(TAG_CACHE_SLUG)
	for _, entry := range cache {
		_, isCached := actualCache[entry.NormalizedName]
		if !isCached {
			actualCache[entry.NormalizedName] = entry
		}
	}

	t.tagCache.Set(TAG_CACHE_SLUG, actualCache, TAG_CACHE_DURATION)

	return tags
}

const TAG_CACHE_DURATION = time.Hour * 24
const TAG_CACHE_SLUG = "Tags"
