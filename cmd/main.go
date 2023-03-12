package main

import (
	"passvault/pkg/app"

	"github.com/gorilla/mux"
)

func main() {
	app := app.App{
		AppRouter: mux.NewRouter(),
		AppConfig: *app.GetAppConfig(),
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
