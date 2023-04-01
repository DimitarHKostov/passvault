package singleton

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/generator"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
)

func GetCookieManager() *cookie.CookieManager {
	return cookie.Get()
}

func GetJwtManager() *jwt.JWTManager {
	return jwt.Get()
}

func GetDatabaseManager() *database.DatabaseManager {
	return database.Get()
}

func GetCryptManager() *crypt.CryptManager {
	return crypt.Get()
}

func GetLogManager() *log.LogManager {
	return log.Get()
}

func GetPayloadGenerator() *generator.PayloadGenerator {
	return generator.Get()
}
