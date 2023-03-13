package api

import (
	"log"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/types"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		for _, cookie := range cookies {
			if cookie != nil && cookie.Name == types.CookieName {
				if cookie.Value == "" {
					log.Println("empty cookie")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

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

				if sessionManager.Get() != cookie.Value {
					log.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}