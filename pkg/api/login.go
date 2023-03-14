package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		log.Println("empty body")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var credentials types.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.LoginValidation{PasswordToValidate: []byte(credentials.Password)}
	if err := validation.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	credentials.Password = string(hashManager.Hash(credentials.Password))

	cookie, err := cookieManager.Produce(types.CookieName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
