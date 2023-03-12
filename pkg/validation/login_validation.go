package validation

import (
	"errors"
)

const (
	invalidCredentialsErrorMessage = "invalid credentials"
)

var (
	actualPassword          = []byte{6, 45, 163, 29, 139, 144, 186, 232, 229, 145, 229, 154, 179, 49, 76, 122, 116, 55, 181, 62, 12, 220, 249, 39, 68, 205, 220, 215, 72, 152, 186, 168}
	invalidCredentialsError = errors.New(invalidCredentialsErrorMessage)
)

type LoginValidation struct {
	PasswordToValidate []byte
}

func (l *LoginValidation) Validate() error {
	for i, e := range l.PasswordToValidate {
		if e != actualPassword[i] {
			return invalidCredentialsError
		}
	}

	return nil
}
