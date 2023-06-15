package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"passvault/pkg/cookie"
	"passvault/pkg/crypt"
	"passvault/pkg/database"
	"passvault/pkg/log"
	"passvault/pkg/middleware"
	"passvault/pkg/operation"
	"passvault/pkg/types"
	"passvault/pkg/validation"

	"github.com/gorilla/mux"
)

const (
	basePathTemplate           = "/api/%s"
)

var (
	basePath *string
	app      *App
)

type App struct {
	appRouter       *mux.Router
	appConfig       *AppConfig
	logManager      *log.LogManagerInterface
	environment     *types.Environment
	databaseManager *database.DatabaseManagerInterface
	cryptManager    *crypt.CryptManagerInterface
	cookieManager   *cookie.CookieManagerInterface
}

// todo refactor at some point(pattern WITH probably)
func NewApp(appRouter *mux.Router, appConfig *AppConfig, logManager *log.LogManagerInterface, environment *types.Environment, databaseManager *database.DatabaseManagerInterface, cryptManager *crypt.CryptManagerInterface, cookieManager *cookie.CookieManagerInterface) *App {
	app := &App{
		appRouter:       appRouter,
		appConfig:       appConfig,
		logManager:      logManager,
		environment:     environment,
		databaseManager: databaseManager,
		cryptManager:    cryptManager,
		cookieManager:   cookieManager}

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
	secretKey := a.environment.JWTSecretKey

	//todo refactor at some point
	a.addEndpoint(a.constructPath(operation.Login), a.login, http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Save), middleware.Middleware(http.HandlerFunc(a.save), secretKey), http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Retrieve), middleware.Middleware(http.HandlerFunc(a.retrieve), secretKey), http.MethodGet)
	a.addEndpoint(a.constructPath(operation.Update), middleware.Middleware(http.HandlerFunc(a.update), secretKey), http.MethodPut)
}

func (a *App) Run() error {
	(*a.logManager).LogInfo(initMessage)

	a.registerEndpoints()
	return http.ListenAndServe(a.appConfig.appPort, a.appRouter)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.logManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var credentials types.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.LoginValidation{PasswordToValidate: []byte(credentials.Password)}
	if err := validation.Validate(); err != nil {
		(*a.logManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := (*a.cookieManager).ProduceCookie()
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (a *App) save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.logManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		(*a.logManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.databaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found {
		(*a.logManager).LogDebug(domainAlreadyExistsMessage)
		w.WriteHeader(http.StatusConflict)
		return
	}

	encryptedPassword, err := (*a.cryptManager).Encrypt(entry.Password)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = (*a.databaseManager).Save(entry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) retrieve(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.logManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.DomainValidation{DomainToValidate: entry.Domain}
	if err := validation.Validate(); err != nil {
		(*a.logManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.databaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		(*a.logManager).LogDebug(domainDoesNotExistMessage)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	queriedEntry, err := (*a.databaseManager).Get(entry.Domain)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := (*a.cryptManager).Decrypt(queriedEntry.Password)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queriedEntry.Password = *decryptedPassword

	jsonBytes, err := json.Marshal(&queriedEntry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		return
	}

	w.Write(jsonBytes)
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.logManager).LogError(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.databaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		(*a.logManager).LogDebug(domainDoesNotExistMessage)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encryptedPassword, err := (*a.cryptManager).Encrypt(entry.Password)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = (*a.databaseManager).Update(entry)
	if err != nil {
		(*a.logManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
