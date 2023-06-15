package cookie

import (
	"fmt"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
	"passvault/pkg/types"
	"time"
)

const (
	expirationTime = 10 * time.Minute
)

type CookieManager struct {
	jwtManager jwt.JWTManagerInterface
	logManager log.LogManagerInterface
}

func NewCookieManager(jwtManager jwt.JWTManagerInterface, logManager log.LogManagerInterface) *CookieManager {
	cookieManager := &CookieManager{
		jwtManager: jwtManager,
		logManager: logManager,
	}

	return cookieManager
}

func (c *CookieManager) ProduceCookie() (*http.Cookie, error) {
	token, err := c.jwtManager.GenerateToken(expirationTime)
	if err != nil {
		c.logManager.LogError(fmt.Sprintf(errorTokenGenerationMessage, err))
		return nil, err
	}

	cookie := http.Cookie{Name: types.CookieName, Value: token, Expires: time.Now().Add(expirationTime), HttpOnly: types.CookieHttpOnly}

	c.logManager.LogDebug(successfulCookieCreationMessage)

	return &cookie, nil
}
