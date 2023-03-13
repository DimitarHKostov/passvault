package cookie

import (
	"errors"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/types"
	"time"
)

const (
	expirationTime = 30 * time.Minute
)

var (
	cookieManager *CookieManager
)

type CookieManager struct {
	JWTManager *jwt.JWTManager
}

func Get() *CookieManager {
	if cookieManager == nil {
		cookieManager = &CookieManager{
			JWTManager: jwt.Get(),
		}
	}

	return cookieManager
}

func (c *CookieManager) Produce(name string, credentials types.Credentials) (*http.Cookie, error) {
	token, err := c.JWTManager.GenerateToken(expirationTime)
	if err != nil {
		return nil, errors.New("error occurred while creating token")
	}

	cookie := http.Cookie{Name: types.CookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: true}

	return &cookie, nil
}
