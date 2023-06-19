package main

import (
	"os"
	"passvault/pkg/app"
	"passvault/pkg/log"
	"passvault/pkg/singleton"
	"passvault/pkg/types"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	ll "log"
)

const (
	envJwtSecretKey     = "JWT_SECRET_KEY"
	envCrypterSecretKey = "CRYPTER_SECRET_KEY"
	envDbHost           = "DB_HOST"
	envDbPort           = "DB_PORT"
	envDbUsername       = "DB_USERNAME"
	envDbPassword       = "DB_PASSWORD"
	envDbName           = "DB_NAME"
	envLogLevel         = "LOG_LEVEL"
	envVaultPassword    = "VAULT_PASSWORD"
)

func main() {
	app := initApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initApp() *app.App {
	app := app.NewApp(withLogManager, withAppRouter, withCookieManager, withCryptManager, withDatabaseManager, withMiddleware, withVaultPasswordHash)

	return app
}

func withVaultPasswordHash(opts *app.AppOpts) {
	env := getEnvironmentVariables()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(env.VaultPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	ll.Println("e", env.VaultPassword)

	opts.VaultPassword = hashedPassword
}

func withMiddleware(opts *app.AppOpts) {
	env := getEnvironmentVariables()
	middleware := singleton.GetMiddleware(env)

	opts.Middleware = middleware
}

func withCookieManager(opts *app.AppOpts) {
	env := getEnvironmentVariables()
	cookieManager := singleton.GetCookieManager(env)

	opts.CookieManager = cookieManager
}

func withDatabaseManager(opts *app.AppOpts) {
	env := getEnvironmentVariables()
	databaseManager := singleton.GetDatabaseManager(env)

	opts.DatabaseManager = databaseManager
}

func withLogManager(opts *app.AppOpts) {
	env := getEnvironmentVariables()

	withLogLevel := func(logOpts *log.LogOpts) {
		logOpts.Level = env.LogLevel
	}

	logManager := singleton.GetLogManager(withLogLevel)

	opts.LogManager = logManager
}

func withAppRouter(opts *app.AppOpts) {
	appRouter := mux.NewRouter()

	opts.AppRouter = appRouter
}

func withCryptManager(opts *app.AppOpts) {
	env := getEnvironmentVariables()
	cryptManager := singleton.GetCryptManager(env)

	opts.CryptManager = cryptManager
}

func getEnvironmentVariables() *types.Environment {
	return &types.Environment{
		JWTSecretKey:     os.Getenv(envJwtSecretKey),
		CrypterSecretKey: os.Getenv(envCrypterSecretKey),
		DbHost:           os.Getenv(envDbHost),
		DbPort:           os.Getenv(envDbPort),
		DbUsername:       os.Getenv(envDbUsername),
		DbPassword:       os.Getenv(envDbPassword),
		DbName:           os.Getenv(envDbName),
		LogLevel:         getEnvVarOrDefault(os.Getenv(envLogLevel), types.DefaultLogLevel),
		VaultPassword:    []byte(os.Getenv(envVaultPassword)),
	}
}

func getEnvVarOrDefault(envVar, defaultValue string) string {
	if envVar != "" {
		return envVar
	}

	return defaultValue
}
