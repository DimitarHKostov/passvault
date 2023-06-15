package middleware

import (
	"net/http"
	"passvault/pkg/jwt"

	//"passvault/pkg/log"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
)

const (
	emptyCookieValueMessage = "cookie value not provided"
)

func Middleware(next http.HandlerFunc, secretKey string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logManager := singleton.GetLogManager()
		cookies := r.Cookies()

		for _, cookie := range cookies {
			if cookie != nil && cookie.Name == types.CookieName {
				if cookie.Value == "" {
					logManager.LogError(emptyCookieValueMessage)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				jwtManager := singleton.GetJwtManager(secretKey)

				_, err := jwtManager.VerifyToken(cookie.Value)
				if err == jwt.InvalidTokenError || err == jwt.ExpiredTokenError {
					logManager.LogError(err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if err != nil {
					logManager.LogError(internetServerErrorMessage)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				logManager.LogDebug(successfulMiddlewareCheckMessage)
				next.ServeHTTP(w, r)
				return
			}
		}

		logManager.LogDebug(cookieNotProvidedMessage)
		w.WriteHeader(http.StatusUnauthorized)
	})
}
