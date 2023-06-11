package main

import (
	"passvault/pkg/app"
	"passvault/pkg/singleton"

	"github.com/gorilla/mux"
)

func run() {
	app := app.NewApp(mux.NewRouter(), *app.GetAppConfig(), singleton.GetLogManager())

	if err := app.Run(); err != nil {
		//todo log
		panic(err)
	}
}

func main() {
	run()
}
