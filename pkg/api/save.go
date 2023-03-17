package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

func Save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		log.Println("empty body")
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

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	databaseManager := singleton.GetDatabaseManager()

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found {
		w.WriteHeader(http.StatusConflict)
		return
	}

	cryptManager := singleton.GetCryptManager()

	encryptedPassword, err := cryptManager.Encrypt(entry.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = encryptedPassword

	err = databaseManager.Save(entry)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
