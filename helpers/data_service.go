package helpers

import (
	"fmt"
	"main/data"
)

func GetEntity[T any](id int) (T, error) {
	var target T

	result := data.DB.First(&target, id)
	if result.RowsAffected != 1 {
		if result.Error != nil {
			return target, result.Error
		}

		return target, fmt.Errorf("no items of type '%s' found", GetStructNameAsString(target))
	}

	return target, nil
}

func GetEntityBySpotifyId[T any](spotifyId string) (T, error) {
	var target T

	result := data.DB.Where("spotify_id = ?", spotifyId).First(&target)
	if result.Error != nil {
		return target, result.Error
	}

	return target, nil
}

func GetEntitiesBySpotifyId[T any](spotifyId string) ([]T, error) {
	var target []T

	result := data.DB.Where("spotify_id = ?", spotifyId).Find(&target)
	if result.Error != nil {
		return target, result.Error
	}

	return target, nil
}
