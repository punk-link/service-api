package helpers

import "reflect"

func GetStructNameAsString[T any](target T) string {
	return reflect.TypeOf(target).Name()
}
