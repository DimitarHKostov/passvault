package app

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/log"
	"passvault/pkg/middleware"
	"passvault/pkg/types"

	"github.com/gorilla/mux"
)

type AppOptFunc func(*AppOpts)

type AppOpts struct {
	AppRouter       *mux.Router
	LogManager      log.LogManagerInterface
	Environment     types.Environment
	DatabaseManager database.DatabaseManagerInterface
	CryptManager    crypt.CryptManagerInterface
	CookieManager   cookie.CookieManagerInterface
	Middleware      middleware.MiddlewareInterface
}

func defaultAppOpts() AppOpts {
	return AppOpts{
		AppRouter: mux.NewRouter(),
	}
}
