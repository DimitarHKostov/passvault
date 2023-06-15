package main

import (
	"passvault/pkg/app"
	"passvault/pkg/database"
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
	app := app.NewApp(withLogManager, withAppRouter, withCookieManager, withCryptManager, withDatabaseManager, withEnvironment)

	return app
}

func withEnvironment(opts *app.AppOpts) {
	envVariables := getEnvironmentVariables()

	opts.Environment = envVariables
}

func withCookieManager(opts *app.AppOpts) {
	envVariables := getEnvironmentVariables()
	cookieManager := singleton.GetCookieManager(envVariables.JWTSecretKey)

	opts.CookieManager = &cookieManager
}

func withDatabaseManager(opts *app.AppOpts) {
	envVariables := getEnvironmentVariables()
	databaseConfig := database.NewDatabaseConfig(envVariables.Host, envVariables.Port, envVariables.Username, envVariables.Password, envVariables.DatabaseName)
	databaseManager := singleton.GetDatabaseManager(*databaseConfig)

	opts.DatabaseManager = &databaseManager
}

func withLogManager(opts *app.AppOpts) {
	logManager := singleton.GetLogManager()

	opts.LogManager = &logManager
}

func withAppRouter(opts *app.AppOpts) {
	appRouter := mux.NewRouter()

	opts.AppRouter = appRouter
}

func withCryptManager(opts *app.AppOpts) {
	envVariables := getEnvironmentVariables()
	cryptManager := singleton.GetCryptManager([]byte(envVariables.CrypterSecretKey))

	opts.CryptManager = &cryptManager
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
	}
}
