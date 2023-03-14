package api

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/jwt"
)

var (
	cookieManager   = cookie.Get()
	jwtManager      = jwt.Get()
	databaseManager = database.Get()
	cryptManager    = crypt.Get()
)

const (
	emptyBodyMessage        = "empty body"
	emptyCookieValueMessage = "cookie value not provided"
)
