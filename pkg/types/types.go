package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Credentials struct {
	Password string `json:"password"`
}

type Payload struct {
	Uuid      uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}

const CookieName = "passvault-cookie"
