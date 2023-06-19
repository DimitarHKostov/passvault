package types

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	CookieName       = "passvault-cookie"
	EmptyBodyMessage = "empty body"
	CookieHttpOnly   = true
	DefaultLogLevel  = "info"
)

type Credentials struct {
	Password string `json:"password"`
}

type Payload struct {
	Uuid      uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Entry struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		errorMessage := "token has expired"

		return errors.New(errorMessage)
	}

	return nil
}

type Environment struct {
	LogLevel         string
	JWTSecretKey     string
	CrypterSecretKey string
	DbHost           string
	DbPort           string
	DbUsername       string
	DbPassword       string
	DbName           string
	VaultPassword    []byte
}
