package cookie

import (
	"errors"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/types"
	"time"
)

const (
	cookieName = "passvault-cookie"
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

func (c *CookieManager) Produce(credentials types.Credentials) (*http.Cookie, error) {
	token, err := c.JWTManager.GenerateToken(5 * time.Minute)
	if err != nil {
		return nil, errors.New("error occurred while creating token")
	}

	cookie := http.Cookie{Name: cookieName, Value: token, Expires: time.Now().Add(5 * time.Minute), HttpOnly: true}
	
	return &cookie, nil
}
