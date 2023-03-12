package hash

import (
	"crypto/sha256"
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
	newHasher := sha256.New()
	newHasher.Write([]byte(str))
	hashed := newHasher.Sum(nil)

	return hashed
}
