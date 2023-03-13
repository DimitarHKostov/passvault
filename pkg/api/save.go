package api

import (
	"net/http"
)

func Save(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
