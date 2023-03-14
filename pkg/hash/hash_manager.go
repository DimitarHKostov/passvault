package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type HashManager struct{}

var (
	hashManager *HashManager
)

func Get() *HashManager {
	if hashManager == nil {
		hashManager = &HashManager{}
	}

	return hashManager
}

func (h *HashManager) Hash(str string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return hashedPassword
}
