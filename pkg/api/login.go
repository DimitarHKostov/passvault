package api

import (
	"encoding/json"
	"io"
	"net/http"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Login(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.Logger.Debug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var credentials types.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		logManager.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.LoginValidation{PasswordToValidate: []byte(credentials.Password), LogManager: logManager}
	if err := validation.Validate(); err != nil {
		logManager.Logger.Debug(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieManager := singleton.GetCookieManager()

	cookie, err := cookieManager.ProduceCookie()
	if err != nil {
		logManager.Logger.Debug(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
