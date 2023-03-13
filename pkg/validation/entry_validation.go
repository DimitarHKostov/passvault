package validation

import (
	"errors"
	"passvault/pkg/types"
)

const (
	argumentsNotSufficientErrorMessage = "args not sufficient"
)

var (
	argumentsNotSufficientError = errors.New(argumentsNotSufficientErrorMessage)
)

type EntryValidation struct {
	EntryToValidate types.Entry
}

func (l *EntryValidation) Validate() error {
	if l.EntryToValidate.Domain == "" || l.EntryToValidate.Username == "" || l.EntryToValidate.Password == "" {
		return argumentsNotSufficientError
	}

	return nil
}
