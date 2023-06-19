package main

import (
	"passvault/pkg/app"
	"passvault/pkg/log"
	"passvault/pkg/singleton"
	"passvault/pkg/types"

	"github.com/gorilla/mux"
)

func main() {
	app := initApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initApp() *app.App {
	app := app.NewApp(withLogManager, withAppRouter, withCookieManager, withCryptManager, withDatabaseManager, withMiddleware)

	return app
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
		JWTSecretKey:     "asdasasdasasdasasdasasdasaa",
		CrypterSecretKey: "this is secret key enough 32 bit",
		Host:             "localhost",
		Port:             "3306",
		Username:         "root",
		Password:         "password",
		DatabaseName:     "db",
		LogLevel:         "debug",
	}
}
