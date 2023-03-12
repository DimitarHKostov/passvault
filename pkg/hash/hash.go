package hash

import (
	"crypto/sha256"
)

type Hasher struct{}

var (
	hasher *Hasher
)

func Get() *Hasher {
	if hasher == nil {
		hasher = &Hasher{}
	}

	return hasher
}

func (h *Hasher) Hash(str string) []byte {
	newHasher := sha256.New()
	newHasher.Write([]byte(str))
	hashed := newHasher.Sum(nil)

	return hashed
}
