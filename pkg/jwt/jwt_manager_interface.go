package jwt

import (
	"passvault/pkg/types"
	"time"
)

type JWTManagerInterface interface {
	GenerateToken(duration time.Duration) (string, error)
	VerifyToken(token string) (*types.Payload, error)
}
