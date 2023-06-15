package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"passvault/pkg/middleware"
	"passvault/pkg/operation"
	"passvault/pkg/types"
	"passvault/pkg/validation"
)

const (
	basePathTemplate = "/api/%s"
)

var (
	basePath *string
	app      *App
)

type App struct {
	appConfig *AppConfig
	AppOpts
}

func NewApp(opts ...AppOptFunc) *App {
	appOpts := defaultAppOpts()

	for _, fn := range opts {
		fn(&appOpts)
	}

	app := &App{AppOpts: appOpts, appConfig: NewAppConfig()}

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
	a.AppOpts.AppRouter.Path(path).HandlerFunc(handlerFunc).Methods(methods...)
}

func (a *App) registerEndpoints() {
	secretKey := a.AppOpts.Environment.JWTSecretKey

	//todo refactor at some point
	a.addEndpoint(a.constructPath(operation.Login), a.login, http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Save), middleware.Middleware(http.HandlerFunc(a.save), secretKey), http.MethodPost)
	a.addEndpoint(a.constructPath(operation.Retrieve), middleware.Middleware(http.HandlerFunc(a.retrieve), secretKey), http.MethodGet)
	a.addEndpoint(a.constructPath(operation.Update), middleware.Middleware(http.HandlerFunc(a.update), secretKey), http.MethodPut)
}

func (a *App) Run() error {
	(*a.AppOpts.LogManager).LogInfo(initMessage)

	a.registerEndpoints()
	return http.ListenAndServe(a.appConfig.appPort, a.AppOpts.AppRouter)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.AppOpts.LogManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var credentials types.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.LoginValidation{PasswordToValidate: []byte(credentials.Password)}
	if err := validation.Validate(); err != nil {
		(*a.AppOpts.LogManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := (*a.AppOpts.CookieManager).ProduceCookie()
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (a *App) save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.AppOpts.LogManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		(*a.AppOpts.LogManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.AppOpts.DatabaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found {
		(*a.AppOpts.LogManager).LogDebug(domainAlreadyExistsMessage)
		w.WriteHeader(http.StatusConflict)
		return
	}

	encryptedPassword, err := (*a.AppOpts.CryptManager).Encrypt(entry.Password)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = (*a.AppOpts.DatabaseManager).Save(entry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) retrieve(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.AppOpts.LogManager).LogDebug(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.DomainValidation{DomainToValidate: entry.Domain}
	if err := validation.Validate(); err != nil {
		(*a.AppOpts.LogManager).LogDebug(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.AppOpts.DatabaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		(*a.AppOpts.LogManager).LogDebug(domainDoesNotExistMessage)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	queriedEntry, err := (*a.AppOpts.DatabaseManager).Get(entry.Domain)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := (*a.AppOpts.CryptManager).Decrypt(queriedEntry.Password)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queriedEntry.Password = *decryptedPassword

	jsonBytes, err := json.Marshal(&queriedEntry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		return
	}

	w.Write(jsonBytes)
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		(*a.AppOpts.LogManager).LogError(types.EmptyBodyMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var entry types.Entry
	err = json.Unmarshal(body, &entry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validation := validation.EntryValidation{EntryToValidate: entry}
	if err := validation.Validate(); err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	found, err := (*a.AppOpts.DatabaseManager).Contains(entry.Domain)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		(*a.AppOpts.LogManager).LogDebug(domainDoesNotExistMessage)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encryptedPassword, err := (*a.AppOpts.CryptManager).Encrypt(entry.Password)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	entry.Password = *encryptedPassword

	err = (*a.AppOpts.DatabaseManager).Update(entry)
	if err != nil {
		(*a.AppOpts.LogManager).LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
