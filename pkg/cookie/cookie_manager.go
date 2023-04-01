package cookie

import (
	"errors"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
	"passvault/pkg/types"
	"time"
)

const (
	expirationTime = 10 * time.Minute
)

var (
	cookieManager *CookieManager
)

type CookieManager struct {
	JWTManager *jwt.JWTManager
	LogManager log.LogManagerInterface
}

func Get() *CookieManager {
	if cookieManager == nil {
		cookieManager = &CookieManager{
			JWTManager: jwt.Get(),
			LogManager: log.Get(),
		}
	}

	return cookieManager
}

func (c *CookieManager) ProduceCookie() (*http.Cookie, error) {
	token, err := c.JWTManager.GenerateToken(expirationTime)
	if err != nil {
		//todo log
		errorMessage := "error occurred while creating token"
		return nil, errors.New(errorMessage)
	}

	//todo log
	cookie := http.Cookie{Name: types.CookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: true}

	return &cookie, nil
}
