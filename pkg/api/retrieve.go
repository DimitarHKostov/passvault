package api

import (
	"encoding/json"
	"io"
	"net/http"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Retrieve(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogDebug(types.EmptyBodyMessage)
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

	validation := validation.DomainValidation{DomainToValidate: entry.Domain}
	if err := validation.Validate(); err != nil {
		logManager.LogDebug(err.Error())
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

	queriedEntry, err := databaseManager.Get(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cryptManager := singleton.GetCryptManager()

	decryptedPassword, err := cryptManager.Decrypt(queriedEntry.Password)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queriedEntry.Password = *decryptedPassword

	jsonBytes, err := json.Marshal(&queriedEntry)
	if err != nil {
		logManager.LogError(err.Error())
		return
	}

	w.Write(jsonBytes)
}
