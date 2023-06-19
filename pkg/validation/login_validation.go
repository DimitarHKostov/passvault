package validation

import (
	"golang.org/x/crypto/bcrypt"
)

// var (
// 	actualPassword = []byte{36, 50, 97, 36, 49, 48, 36, 112, 73, 116, 72, 106, 97, 80, 57, 52, 49, 106, 57, 100, 71, 74, 87, 76, 51, 72, 109, 49, 101, 53, 57, 113, 117, 107, 107, 113, 119, 121, 86, 51, 52, 112, 117, 72, 121, 51, 100, 50, 113, 83, 82, 119, 86, 54, 80, 67, 49, 112, 115, 79}
// )

type LoginValidation struct {
	ActualPassword     []byte
	PasswordToValidate []byte
}

func (lv *LoginValidation) Validate() error {
	return bcrypt.CompareHashAndPassword(lv.ActualPassword, lv.PasswordToValidate)
}
