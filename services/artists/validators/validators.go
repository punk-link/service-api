package validators

import (
	"errors"
	data "main/data/artists"
)

func CurrentDbArtistBelongsToLabel(err error, dbArtist data.Artist, targetLabelId int) error {
	if err != nil {
		return err
	}

	if dbArtist.LabelId != targetLabelId {
		err = errors.New("artist already added to another label")
	}

	return err
}
