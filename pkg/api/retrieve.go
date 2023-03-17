package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Retrieve(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		log.Println(emptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.DomainValidation{DomainToValidate: entry.Domain}
	if err := validation.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	queriedEntry, err := databaseManager.Get(entry.Domain)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := cryptManager.Decrypt(queriedEntry.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queriedEntry.Password = decryptedPassword

	jsonBytes, err := json.Marshal(&queriedEntry)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(jsonBytes)
}
