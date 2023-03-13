package api

import (
	"passvault/pkg/cookie"
	"passvault/pkg/database"
	"passvault/pkg/hash"
	"passvault/pkg/jwt"
	"passvault/pkg/session"
)

var (
	cookieManager   = cookie.Get()
	hashManager     = hash.Get()
	jwtManager      = jwt.Get()
	sessionManager  = session.Get()
	databaseManager = database.Get()
)
