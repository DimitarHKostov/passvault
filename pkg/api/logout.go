package api

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionManager.InvalidateSession()
	w.WriteHeader(http.StatusOK)
}
