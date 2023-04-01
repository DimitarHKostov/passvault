package cookie

import "net/http"

type CookieManagerInterface interface {
	ProduceCookie() (*http.Cookie, error)
}
