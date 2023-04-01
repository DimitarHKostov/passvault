package singleton

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/generator"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
)

func GetCookieManager() cookie.CookieManagerInterface {
	return cookie.Get()
}

func GetJwtManager() jwt.JWTManagerInterface {
	return jwt.Get()
}

func GetDatabaseManager() database.DatabaseManagerInterface {
	return database.Get()
}

func GetCryptManager() crypt.CryptManagerInterface {
	return crypt.Get()
}

func GetLogManager() log.LogManagerInterface {
	return log.Get()
}

func GetPayloadGenerator() generator.PayloadGeneratorInterface {
	return generator.Get()
}
