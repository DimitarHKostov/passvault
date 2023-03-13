package api

import (
	"net/http"
)

func Retrieve(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
