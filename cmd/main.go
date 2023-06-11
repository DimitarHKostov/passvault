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
		SecretKey: "asdasasdasasdasasdasasdasaa",
	}
}
