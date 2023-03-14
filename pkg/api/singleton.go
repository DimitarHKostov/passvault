package api

import (
	"passvault/pkg/cookie"
	"passvault/pkg/database"
	"passvault/pkg/hash"
	"passvault/pkg/jwt"
)

var (
	cookieManager = cookie.Get()
	hashManager   = hash.Get()
	jwtManager    = jwt.Get()
	databaseManager = database.Get()
)
