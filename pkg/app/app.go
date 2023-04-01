package app

import (
	"fmt"
	"net/http"
	"passvault/pkg/handler_func_factory"
	"passvault/pkg/operation"

	"github.com/gorilla/mux"
)

const (
	basePathTemplate = "/api/%s"
)

var (
	basePath *string
)

type App struct {
	AppRouter *mux.Router
	AppConfig AppConfig
}

func (a *App) constructPath(operation operation.Operation) string {
	if basePath == nil {
		basePath = new(string)
		*basePath = fmt.Sprintf(basePathTemplate, a.AppConfig.AppVersion)
	}

	return *basePath + fmt.Sprintf("/%s", operation.String())
}

func (a *App) addEndpoint(path string, handlerFunc func(http.ResponseWriter, *http.Request), methods ...string) {
	a.AppRouter.PathPrefix(path).HandlerFunc(handlerFunc).Methods(methods...)
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

	return http.ListenAndServe(a.AppConfig.AppPort, a.AppRouter)
}
