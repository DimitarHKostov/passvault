package app

import (
	"fmt"
	"net/http"
	"passvault/pkg/handler_func_factory"
	"passvault/pkg/log"
	"passvault/pkg/operation"

	"github.com/gorilla/mux"
)

const (
	basePathTemplate = "/api/%s"
)

var (
	basePath *string
	app      *App
)

type App struct {
	appRouter  *mux.Router
	appConfig  AppConfig
	logManager log.LogManagerInterface
}

func NewApp(appRouter *mux.Router, appConfig AppConfig, logManager log.LogManagerInterface) *App {
	if app == nil {
		app = &App{appRouter: appRouter, appConfig: appConfig, logManager: logManager}
	}

	return app
}

func (a *App) constructPath(operation operation.Operation) string {
	if basePath == nil {
		basePath = new(string)
		*basePath = fmt.Sprintf(basePathTemplate, a.appConfig.AppVersion)
	}

	return *basePath + fmt.Sprintf("/%s", operation.String())
}

func (a *App) addEndpoint(path string, handlerFunc func(http.ResponseWriter, *http.Request), methods ...string) {
	a.appRouter.Path(path).HandlerFunc(handlerFunc).Methods(methods...)
}

func (a *App) registerEndpoints() {
	handlerFuncFactory := handler_func_factory.Get()

	loginHandlerFunc := handlerFuncFactory.Produce(operation.Login)
	saveHandlerFunc := handlerFuncFactory.Produce(operation.Save)
	retrieveHandlerFunc := handlerFuncFactory.Produce(operation.Retrieve)
	updateHandlerFunc := handlerFuncFactory.Produce(operation.Update)

	a.addEndpoint(a.constructPath(operation.Login), loginHandlerFunc, http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Save), saveHandlerFunc, http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Retrieve), retrieveHandlerFunc, http.MethodGet)
	a.addEndpoint(a.constructPath(operation.Update), updateHandlerFunc, http.MethodPut)
}

func (a *App) Run() error {
	a.registerEndpoints()

	//todo log

	return http.ListenAndServe(a.appConfig.AppPort, a.appRouter)
}
