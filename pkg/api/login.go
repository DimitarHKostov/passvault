package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"passvault/pkg/hash"
	"passvault/pkg/jwt"
	"passvault/pkg/types"
	"passvault/pkg/validation"
	"time"
)

var (
	hasher     = hash.Get()
	jwtManager = jwt.Get()
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

	token, err := jwtManager.GenerateToken(5 * time.Minute)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "passvault-cookie", Value: token, Expires: time.Now().Add(5 * time.Minute), HttpOnly: true}
	http.SetCookie(w, &cookie)

	w.Write([]byte("success"))
}
