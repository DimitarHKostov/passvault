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
	minSecretKeySize         = 32
	invalidTokenErrorMessage = "token is invalid"
	expiredTokenErrorMessage = "token has expired"
	secretKey                = "asdasasdasasdasasdasasdasaa"
)

var (
	InvalidTokenError = errors.New(invalidTokenErrorMessage)
	ExpiredTokenError = errors.New(expiredTokenErrorMessage)
	jwtManager        *JWTManager
)

type JWTManager struct {
	PayloadGenerator generator.PayloadGeneratorInterface
	SecretKey        string
	LogManager       *log.LogManager
}

func Get() *JWTManager {
	if jwtManager == nil {
		jwtManager = &JWTManager{
			PayloadGenerator: generator.Get(),
			SecretKey:        secretKey,
			LogManager:       log.Get(),
		}
	}

	return jwtManager
}

func (jwtm *JWTManager) GenerateToken(duration time.Duration) (string, error) {
	payload, err := jwtm.PayloadGenerator.GeneratePayload(duration)
	if err != nil {
		//todo log
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//todo log
	return jwtToken.SignedString([]byte(jwtm.SecretKey))
}

func (jwtm *JWTManager) VerifyToken(token string) (*types.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			//todo log
			return nil, InvalidTokenError
		}
		//todo log
		return []byte(jwtm.SecretKey), nil
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
