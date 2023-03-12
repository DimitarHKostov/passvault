package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"passvault/pkg/cookie"
	"passvault/pkg/hash"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

var (
	cookieManager = cookie.Get()
	hasher        = hash.Get()
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
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

	hashedPassword := hasher.Hash(credentials.Password)
	validation := validation.LoginValidation{PasswordToValidate: hashedPassword}
	if err := validation.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := cookieManager.Produce(credentials)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
