package validation

import (
	"errors"
	"passvault/pkg/types"
	"strings"
)

const (
	argumentsNotSufficientErrorMessage = "args not sufficient"
	minPasswordSize                    = 8
	minPasswordSizeErrorMessage        = "password too short, must be at least 8 characters"
	forbiddenCharacters                = "\"'`;"
	forbiddenCharactersErrorMessage    = "password contains forbidden character"
)

var (
	argumentsNotSufficientError = errors.New(argumentsNotSufficientErrorMessage)
	domainValidationSingleton   *DomainValidation
)

type EntryValidation struct {
	EntryToValidate types.Entry
}

func (ev *EntryValidation) Validate() error {
	if err := getDomainValidationInstance().validateDomain(ev.EntryToValidate.Domain); err != nil {
		return err
	}

	if err := ev.validatePassword(ev.EntryToValidate.Password); err != nil {
		return err
	}

	if err := ev.validateUsername(ev.EntryToValidate.Username); err != nil {
		return err
	}

	return nil
}

func (ev *EntryValidation) validatePassword(password string) error {
	return ev.genericValidation(password)
}

func (ev *EntryValidation) validateUsername(username string) error {
	return ev.genericValidation(username)
}

func (ev *EntryValidation) genericValidation(str string) error {
	if len(str) < minPasswordSize {
		return errors.New(minPasswordSizeErrorMessage)
	}

	if strings.ContainsAny(str, forbiddenCharacters) {
		return errors.New(forbiddenCharactersErrorMessage)
	}

	return nil
}

func getDomainValidationInstance() *DomainValidation {
	if domainValidationSingleton == nil {
		domainValidationSingleton = &DomainValidation{}
	}

	return domainValidationSingleton
}
