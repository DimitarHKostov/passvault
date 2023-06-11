package singleton

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/generator"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
)

func GetCookieManager(secretKey string) cookie.CookieManagerInterface {
	return cookie.NewCookieManager(GetJwtManager(secretKey), GetLogManager())
}

func GetJwtManager(secretKey string) jwt.JWTManagerInterface {
	return jwt.NewJwtManager(GetPayloadGenerator(), secretKey, GetLogManager())
}

func GetDatabaseManager() database.DatabaseManagerInterface {
	return database.NewDatabaseManager(GetLogManager())
}

func GetCryptManager() crypt.CryptManagerInterface {
	return crypt.NewCryptManager(GetLogManager())
}

func GetLogManager() log.LogManagerInterface {
	return log.NewLogManager()
}

func GetPayloadGenerator() generator.PayloadGeneratorInterface {
	return generator.NewPayloadGenerator(GetLogManager())
}
