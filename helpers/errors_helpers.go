package helpers

import (
	"fmt"
	"strings"
)

func AccumulateErrors(errs ...error) error {
	messages := make([]string, 0)
	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	if len(messages) != 0 {
		return fmt.Errorf(strings.Join(messages, "; "))
	}

	return nil
}

func CombineErrors(err error, newErr error) error {
	if err == nil {
		return newErr
	}

	if newErr == nil {
		return err
	}

	return fmt.Errorf("%s; %s", err.Error(), newErr.Error())
}
