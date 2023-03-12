package api

import (
	"log"
	"net/http"
	"passvault/pkg/jwt"
	"passvault/pkg/types"
)

func Logout(w http.ResponseWriter, r *http.Request) {
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

			sessionManager.InvalidateSession()
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}
