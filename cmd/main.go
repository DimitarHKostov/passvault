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
	envVariables := getEnvironmentVariables()

	logManager := singleton.GetLogManager()
	appConfig := app.NewAppConfig()
	appRouter := mux.NewRouter()
	databaseConfig := database.NewDatabaseConfig(envVariables.Host, envVariables.Port, envVariables.Username, envVariables.Password, envVariables.DatabaseName)
	databaseManager := singleton.GetDatabaseManager(*databaseConfig)
	cryptManager := singleton.GetCryptManager([]byte(envVariables.CrypterSecretKey))
	cookieManager := singleton.GetCookieManager(envVariables.JWTSecretKey)

	app := app.NewApp(appRouter, appConfig, &logManager, envVariables, &databaseManager, &cryptManager, &cookieManager)

	return app
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
