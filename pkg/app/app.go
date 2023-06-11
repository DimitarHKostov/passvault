package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"passvault/pkg/log"
	"passvault/pkg/middleware"
	"passvault/pkg/operation"
	"passvault/pkg/singleton"
	"passvault/pkg/types"
	"passvault/pkg/validation"

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
	appRouter   *mux.Router
	appConfig   AppConfig
	logManager  log.LogManagerInterface
	environment types.Environment
}

func NewApp(appRouter *mux.Router, appConfig AppConfig, logManager log.LogManagerInterface, environment types.Environment) *App {
	if app == nil {
		app = &App{appRouter: appRouter, appConfig: appConfig, logManager: logManager, environment: environment}
	}

	return app
}

func (a *App) constructPath(operation operation.Operation) string {
	if basePath == nil {
		basePath = new(string)
		*basePath = fmt.Sprintf(basePathTemplate, a.appConfig.appVersion)
	}

	return *basePath + fmt.Sprintf("/%s", operation.String())
}

func (a *App) addEndpoint(path string, handlerFunc func(http.ResponseWriter, *http.Request), methods ...string) {
	a.appRouter.Path(path).HandlerFunc(handlerFunc).Methods(methods...)
}

func (a *App) registerEndpoints() {
	secretKey := a.environment.SecretKey

	//todo refactor at some point
	a.addEndpoint(a.constructPath(operation.Login), a.login, http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Save), middleware.Middleware(http.HandlerFunc(a.save), secretKey), http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Retrieve), middleware.Middleware(http.HandlerFunc(a.retrieve), secretKey), http.MethodGet)
	a.addEndpoint(a.constructPath(operation.Update), middleware.Middleware(http.HandlerFunc(a.update), secretKey), http.MethodPut)
}

func (a *App) Run() error {
	a.registerEndpoints()

	//todo log

	return http.ListenAndServe(a.appConfig.appPort, a.appRouter)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var credentials types.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.LoginValidation{PasswordToValidate: []byte(credentials.Password)}
	if err := validation.Validate(); err != nil {
		logManager.LogDebug(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieManager := singleton.GetCookieManager(a.environment.SecretKey)

	cookie, err := cookieManager.ProduceCookie()
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (a *App) save(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		logManager.LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	databaseManager := singleton.GetDatabaseManager()

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found {
		//todo log
		w.WriteHeader(http.StatusConflict)
		return
	}

	cryptManager := singleton.GetCryptManager()

	encryptedPassword, err := cryptManager.Encrypt(entry.Password)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = databaseManager.Save(entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) retrieve(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.DomainValidation{DomainToValidate: entry.Domain}
	if err := validation.Validate(); err != nil {
		logManager.LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cryptManager := singleton.GetCryptManager()
	databaseManager := singleton.GetDatabaseManager()

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		//todo log
		w.WriteHeader(http.StatusNotFound)
		return
	}

	queriedEntry, err := databaseManager.Get(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := cryptManager.Decrypt(queriedEntry.Password)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queriedEntry.Password = *decryptedPassword

	jsonBytes, err := json.Marshal(&queriedEntry)
	if err != nil {
		logManager.LogError(err.Error())
		return
	}

	w.Write(jsonBytes)
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	logManager := singleton.GetLogManager()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		logManager.LogError(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	databaseManager := singleton.GetDatabaseManager()

	found, err := databaseManager.Contains(entry.Domain)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		//todo log
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cryptManager := singleton.GetCryptManager()

	encryptedPassword, err := cryptManager.Encrypt(entry.Password)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = databaseManager.Update(entry)
	if err != nil {
		logManager.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
