package middleware

import "net/http"

type MiddlewareInterface interface {
	Intercept(next http.HandlerFunc) http.HandlerFunc
}
