package main

import (
	"passvault/pkg/app"
	"passvault/pkg/singleton"
	"passvault/pkg/types"

	"github.com/gorilla/mux"
)

func main() {
	logManager := singleton.GetLogManager()
	appConfig := app.NewAppConfig()
	appRouter := mux.NewRouter()
	envVariables := getEnvironmentVariables()

	app := app.NewApp(appRouter, appConfig, &logManager, envVariables)

	if err := app.Run(); err != nil {
		logManager.LogPanic(err.Error())
	}
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
