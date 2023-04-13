package artists

import (
	"fmt"
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
		insertionErr := t.db.Exec("insert into tags (name, normalized_name, name_tokens) values (?, ?, to_tsvector(?))", tag.Name, tag.NormalizedName, tag.Name).Error
		if insertionErr != nil {
			err = helpers.CombineErrors(err, insertionErr)
		}
	}

	return t.handleError(err)
}

func (t *TagRepositoryService) Get(err error, normalizedNames []string) []artistData.Tag {
	if err != nil {
		return make([]artistData.Tag, 0)
	}

	var tags []artistData.Tag
	err = t.db.Where("normalized_name IN (?)", normalizedNames).
		Find(&tags).
		Error

	t.handleError(err)

	return tags
}

func (t *TagRepositoryService) Search(err error, query string) []artistData.Tag {
	if err != nil {
		return make([]artistData.Tag, 0)
	}

	tags := make([]artistData.Tag, 0)
	var tag artistData.Tag
	rows, err := t.db.Raw(fmt.Sprintf("id, select name, normalized_name from tags where name_tokens @@ to_tsquery('%s')", query)).Rows()
	if err != nil {
		return make([]artistData.Tag, 0)
	}

	defer rows.Close()
	for rows.Next() {
		t.db.ScanRows(rows, &tag)
		tags = append(tags, tag)
	}

	return tags
}

func (t *TagRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}
