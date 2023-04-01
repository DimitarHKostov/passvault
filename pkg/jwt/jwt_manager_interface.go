package jwt

import (
	"passvault/pkg/types"
	"time"
)

type JWTManagerInterface interface {
	GenerateToken(time.Duration) (string, error)
	VerifyToken(string) (*types.Payload, error)
}
