package artists

import (
	artistData "main/data/artists"
	"main/helpers"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type TagRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTagRepository(injector *do.Injector) (TagRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &TagRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *TagRepositoryService) Create(err error, tags *[]artistData.Tag) error {
	if err != nil || len(*tags) == 0 {
		return err
	}

	for _, tag := range *tags {
		insertionErr := t.db.Exec("insert into tags (name, normalized_name, normalized_vector) values (?, ?, to_tsvector(?))", tag.Name, tag.NormalizedName, tag.NormalizedName).Error
		if insertionErr != nil {
			err = helpers.CombineErrors(err, insertionErr)
		}
	}

	return t.handleError(err)
}

func (t *TagRepositoryService) Get(err error, names []string) []artistData.Tag {
	if err != nil {
		return make([]artistData.Tag, 0)
	}

	tags := make([]artistData.Tag, 0)
	//var tag artistData.Tag

	return tags
}

func (t *TagRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}
