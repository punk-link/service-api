package labels

import (
	labelData "main/data/labels"
	"main/helpers"
	labelModels "main/models/labels"
	"main/services/labels/validators"
	spotifyPlatformServices "main/services/platforms/spotify"
	"strings"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type LabelService struct {
	db             *gorm.DB
	logger         logger.Logger
	spotifyService *spotifyPlatformServices.SpotifyService
}

func NewLabelService(injector *do.Injector) (*LabelService, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	spotifyService := do.MustInvoke[*spotifyPlatformServices.SpotifyService](injector)

	return &LabelService{
		db:             db,
		logger:         logger,
		spotifyService: spotifyService,
	}, nil
}

func (t *LabelService) AddLabel(labelName string) (labelModels.Label, error) {
	trimmedName := strings.TrimSpace(labelName)
	err := validators.NameNotEmpty(trimmedName)

	return t.addLabelInternal(err, trimmedName)
}

func (t *LabelService) GetLabel(currentManager labelModels.ManagerContext, id int) (labelModels.Label, error) {
	err := validators.CurrentManagerBelongsToLabel(currentManager, id)
	return t.getLabelWithoutContextCheck(err, id)
}

func (t *LabelService) ModifyLabel(currentManager labelModels.ManagerContext, label labelModels.Label, id int) (labelModels.Label, error) {
	err1 := validators.CurrentManagerBelongsToLabel(currentManager, id)
	err2 := validators.IdConsistsOverRequest(label.Id, id)

	trimmedName := strings.TrimSpace(label.Name)
	err3 := validators.NameNotEmpty(trimmedName)

	return t.modifyLabelInternal(helpers.AccumulateErrors(err1, err2, err3), currentManager, trimmedName)
}

func (t *LabelService) addLabelInternal(err error, labelName string) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
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

func (t *LabelService) getLabelWithoutContextCheck(err error, id int) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
	}

	dbLabel, err := getDbLabel(t.db, t.logger, err, id)
	return labelModels.Label{
		Id:   dbLabel.Id,
		Name: dbLabel.Name,
	}, err
}

func (t *LabelService) modifyLabelInternal(err error, currentManager labelModels.ManagerContext, lebelName string) (labelModels.Label, error) {
	if err != nil {
		return labelModels.Label{}, err
	}

	dbLabel, err := getDbLabel(t.db, t.logger, err, currentManager.LabelId)

	dbLabel.Name = lebelName
	err = updateDbLabel(t.db, t.logger, err, &dbLabel)
	return t.getLabelWithoutContextCheck(err, dbLabel.Id)
}
