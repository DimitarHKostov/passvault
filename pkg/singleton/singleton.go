package singleton

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/generator"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
)

var (
	cookieManager    cookie.CookieManagerInterface
	jwtManager       jwt.JWTManagerInterface
	databaseManager  database.DatabaseManagerInterface
	cryptManager     crypt.CryptManagerInterface
	logManager       log.LogManagerInterface
	payloadGenerator generator.PayloadGeneratorInterface
)

func GetCookieManager(jwtSecretKey string) cookie.CookieManagerInterface {
	if cookieManager == nil {
		cookieManager = cookie.NewCookieManager(GetJwtManager(jwtSecretKey), GetLogManager())
	}

	return cookieManager
}

func GetJwtManager(jwtSecretKey string) jwt.JWTManagerInterface {
	if jwtManager == nil {
		jwtManager = jwt.NewJwtManager(GetPayloadGenerator(), jwtSecretKey, GetLogManager())
	}

	return jwtManager
}

func GetDatabaseManager(databaseConfig database.DatabaseConfig) database.DatabaseManagerInterface {
	if databaseManager == nil {
		databaseManager = database.NewDatabaseManager(GetLogManager(), databaseConfig)
	}

	return databaseManager
}

func GetCryptManager(crypterSecretKey []byte) crypt.CryptManagerInterface {
	if cryptManager == nil {
		cryptManager = crypt.NewCryptManager(GetLogManager(), crypterSecretKey)
	}

	return cryptManager
}

func GetLogManager() log.LogManagerInterface {
	if logManager == nil {
		logManager = log.NewLogManager()
	}

	return logManager
}

func GetPayloadGenerator() generator.PayloadGeneratorInterface {
	if payloadGenerator == nil {
		payloadGenerator = generator.NewPayloadGenerator(GetLogManager())
	}

	return payloadGenerator
}
