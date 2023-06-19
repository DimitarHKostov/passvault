package singleton

import (
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/generator"
	"passvault/pkg/jwt"
	"passvault/pkg/log"
	m "passvault/pkg/middleware"
	"passvault/pkg/types"
)

var (
	cookieManager    cookie.CookieManagerInterface
	jwtManager       jwt.JWTManagerInterface
	databaseManager  database.DatabaseManagerInterface
	cryptManager     crypt.CryptManagerInterface
	logManager       log.LogManagerInterface
	payloadGenerator generator.PayloadGeneratorInterface
	middleware       m.MiddlewareInterface
)

func GetMiddleware(env *types.Environment) m.MiddlewareInterface {
	if middleware == nil {
		middleware = m.NewMiddleware(GetLogManager(), GetJwtManager(env))
	}

	return middleware
}

func GetCookieManager(env *types.Environment) cookie.CookieManagerInterface {
	if cookieManager == nil {
		cookieManager = cookie.NewCookieManager(GetJwtManager(env), GetLogManager())
	}

	return cookieManager
}

func GetJwtManager(env *types.Environment) jwt.JWTManagerInterface {
	if jwtManager == nil {
		jwtManager = jwt.NewJwtManager(GetPayloadGenerator(), env.JWTSecretKey, GetLogManager())
	}

	return jwtManager
}

func GetDatabaseManager(env *types.Environment) database.DatabaseManagerInterface {
	if databaseManager == nil {
		databaseConfig := database.NewDatabaseConfig(env.DbHost, env.DbPassword, env.DbUsername, env.DbPassword, env.DbName)
		databaseManager = database.NewDatabaseManager(GetLogManager(), databaseConfig)
	}

	return databaseManager
}

func GetCryptManager(env *types.Environment) crypt.CryptManagerInterface {
	if cryptManager == nil {
		cryptManager = crypt.NewCryptManager(GetLogManager(), []byte(env.CrypterSecretKey))
	}

	return cryptManager
}

func GetLogManager(opts ...log.LogOptsFn) log.LogManagerInterface {
	if logManager == nil {
		logManager = log.NewLogManager(opts...)
	}

	return logManager
}

func GetPayloadGenerator() generator.PayloadGeneratorInterface {
	if payloadGenerator == nil {
		payloadGenerator = generator.NewPayloadGenerator(GetLogManager())
	}

	return payloadGenerator
}
