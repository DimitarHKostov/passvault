package api

import (
	"encoding/json"
	"io"
	"net/http"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Update(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogError(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry, LogManager: logManager}
	if err := validation.Validate(); err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	databaseManager := singleton.GetDatabaseManager()

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		//todo log
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cryptManager := singleton.GetCryptManager()

	encryptedPassword, err := cryptManager.Encrypt(entry.Password)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = databaseManager.Update(entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
