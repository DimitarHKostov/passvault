package middleware

import (
	"log"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
)

const (
	emptyCookieValueMessage = "cookie value not provided"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		for _, cookie := range cookies {
			if cookie != nil && cookie.Name == types.CookieName {
				if cookie.Value == "" {
					log.Println(emptyCookieValueMessage)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				jwtManager := singleton.GetJwtManager()

				_, err := jwtManager.VerifyToken(cookie.Value)
				if err == jwt.InvalidTokenError || err == jwt.ExpiredTokenError {
					log.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}
