package main

import (
	"passvault/pkg/app"
	"passvault/pkg/singleton"

	"github.com/gorilla/mux"
)

func run() {
	app := app.App{
		AppRouter:  mux.NewRouter(),
		AppConfig:  *app.GetAppConfig(),
		LogManager: singleton.GetLogManager(),
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func main() {
	run()
}
