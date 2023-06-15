package jwt

import (
	"errors"
	"fmt"

	"passvault/pkg/generator"

	"passvault/pkg/log"
	"passvault/pkg/types"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	InvalidTokenError = errors.New(invalidTokenErrorMessage)
	ExpiredTokenError = errors.New(expiredTokenErrorMessage)
)

type JWTManager struct {
	payloadGenerator generator.PayloadGeneratorInterface
	secretKey        string
	logManager       log.LogManagerInterface
}

func NewJwtManager(payloadGenerator generator.PayloadGeneratorInterface, secretKey string, logManager log.LogManagerInterface) *JWTManager {
	jwtManager := &JWTManager{
		payloadGenerator: payloadGenerator,
		secretKey:        secretKey,
		logManager:       logManager,
	}

	return jwtManager
}

func (jwtm *JWTManager) GenerateToken(duration time.Duration) (string, error) {
	payload, err := jwtm.payloadGenerator.GeneratePayload(duration)
	if err != nil {
		jwtm.logManager.LogError(fmt.Sprintf(payloadGenerationFailMessage, err))
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(jwtm.secretKey))
}

func (jwtm *JWTManager) VerifyToken(token string) (*types.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			jwtm.logManager.LogError(invalidTokenErrorMessage)
			return nil, InvalidTokenError
		}

		return []byte(jwtm.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &types.Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredTokenError) {
			jwtm.logManager.LogError(expiredTokenErrorMessage)
			return nil, ExpiredTokenError
		}

		jwtm.logManager.LogError(invalidTokenErrorMessage)

		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*types.Payload)
	if !ok {
		jwtm.logManager.LogError(invalidTokenErrorMessage)
		return nil, InvalidTokenError
	}

	jwtm.logManager.LogDebug(successfulTokenVerification)

	return payload, nil
}
