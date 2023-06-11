package main

import (
	"passvault/pkg/app"
	"passvault/pkg/singleton"
	"passvault/pkg/types"

	"github.com/gorilla/mux"
)

func main() {
	app := app.NewApp(mux.NewRouter(), *app.NewAppConfig(), singleton.GetLogManager(), getEnvironmentVariables())

	if err := app.Run(); err != nil {
		//todo log
		panic(err)
	}
}

func getEnvironmentVariables() types.Environment {
	return types.Environment{
		JWTSecretKey:     "asdasasdasasdasasdasasdasaa",
		CrypterSecretKey: "this is secret key enough 32 bit",
		Host:             "localhost",
		Port:             "3306",
		Username:         "root",
		Password:         "password",
		DatabaseName:     "db",
	}
}
