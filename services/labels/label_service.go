package labels

import (
	labelData "main/data/labels"
	"main/helpers"
	"main/models/labels"
	validator "main/services/labels/validators"
	"main/services/platforms/spotify"
	"strings"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type LabelService struct {
	db             *gorm.DB
	logger         logger.Logger
	spotifyService *spotify.SpotifyService
}

func NewLabelService(injector *do.Injector) (*LabelService, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	spotifyService := do.MustInvoke[*spotify.SpotifyService](injector)

	return &LabelService{
		db:             db,
		logger:         logger,
		spotifyService: spotifyService,
	}, nil
}

func (t *LabelService) AddLabel(labelName string) (labels.Label, error) {
	trimmedName := strings.TrimSpace(labelName)
	err := validator.NameNotEmpty(trimmedName)

	return t.addLabelInternal(err, trimmedName)
}

func (t *LabelService) GetLabel(currentManager labels.ManagerContext, id int) (labels.Label, error) {
	err := validator.CurrentManagerBelongsToLabel(currentManager, id)
	return t.getLabelWithoutContextCheck(err, id)
}

func (t *LabelService) ModifyLabel(currentManager labels.ManagerContext, label labels.Label, id int) (labels.Label, error) {
	err1 := validator.CurrentManagerBelongsToLabel(currentManager, id)
	err2 := validator.IdConsistsOverRequest(label.Id, id)

	trimmedName := strings.TrimSpace(label.Name)
	err3 := validator.NameNotEmpty(trimmedName)

	return t.modifyLabelInternal(helpers.AccumulateErrors(err1, err2, err3), currentManager, trimmedName)
}

func (t *LabelService) addLabelInternal(err error, labelName string) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	now := time.Now().UTC()
	dbLabel := labelData.Label{
		Created: now,
		Name:    labelName,
		Updated: now,
	}

	err = createDbLabel(t.db, t.logger, err, &dbLabel)
	return t.getLabelWithoutContextCheck(err, dbLabel.Id)
}

func (t *LabelService) getLabelWithoutContextCheck(err error, id int) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel, err := getDbLabel(t.db, t.logger, err, id)
	return labels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, err
}

func (t *LabelService) modifyLabelInternal(err error, currentManager labels.ManagerContext, lebelName string) (labels.Label, error) {
	if err != nil {
		return labels.Label{}, err
	}

	dbLabel, err := getDbLabel(t.db, t.logger, err, currentManager.LabelId)

	dbLabel.Name = lebelName
	err = updateDbLabel(t.db, t.logger, err, &dbLabel)
	return t.getLabelWithoutContextCheck(err, dbLabel.Id)
}
