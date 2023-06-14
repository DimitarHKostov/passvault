package jwt

import (
	"errors"

	"passvault/pkg/generator"

	"passvault/pkg/log"
	"passvault/pkg/types"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	invalidTokenErrorMessage = "token is invalid"
	expiredTokenErrorMessage = "token has expired"
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
		//todo log
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//todo log
	return jwtToken.SignedString([]byte(jwtm.secretKey))
}

func (jwtm *JWTManager) VerifyToken(token string) (*types.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			//todo log
			return nil, InvalidTokenError
		}
		//todo log
		return []byte(jwtm.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &types.Payload{}, keyFunc)
	if err != nil {
		//todo log
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredTokenError) {
			//todo log
			return nil, ExpiredTokenError
		}
		//todo log
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*types.Payload)
	if !ok {
		//todo log
		return nil, InvalidTokenError
	}

	//todo log
	return payload, nil
}
