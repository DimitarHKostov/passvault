package middleware

import (
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/log"

	"passvault/pkg/types"
)

const (
	emptyCookieValueMessage = "cookie value not provided"
)

type Middleware struct {
	logManager log.LogManagerInterface
	jwtManager jwt.JWTManagerInterface
}

func NewMiddleware(logManager log.LogManagerInterface, jwtManager jwt.JWTManagerInterface) *Middleware {
	middleware := &Middleware{logManager: logManager, jwtManager: jwtManager}

	return middleware
}

func (m *Middleware) Intercept(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		for _, cookie := range cookies {
			if cookie != nil && cookie.Name == types.CookieName {
				if cookie.Value == "" {
					m.logManager.LogError(emptyCookieValueMessage)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				_, err := m.jwtManager.VerifyToken(cookie.Value)
				if err == jwt.InvalidTokenError || err == jwt.ExpiredTokenError {
					m.logManager.LogError(err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if err != nil {
					m.logManager.LogError(internetServerErrorMessage)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				m.logManager.LogDebug(successfulMiddlewareCheckMessage)
				next.ServeHTTP(w, r)
				return
			}
		}

		m.logManager.LogDebug(cookieNotProvidedMessage)
		w.WriteHeader(http.StatusUnauthorized)
	})
}
