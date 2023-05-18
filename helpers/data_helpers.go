package helpers

func ShouldHandleDbError(err error) bool {
	if err == nil {
		return false
	}

	switch err.Error() {
	case "record not found":
		return false
	}

	return true
}
