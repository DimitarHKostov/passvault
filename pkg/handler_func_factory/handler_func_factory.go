package handler_func_factory

import (
	"net/http"
	"passvault/pkg/api"
	"passvault/pkg/log"
	"passvault/pkg/middleware"
	"passvault/pkg/operation"
)

const (
	invalidOperationMessage = "operatrion is invalid"
)

var (
	handlerFuncFactory *HandlerFuncFactory
)

type HandlerFuncFactory struct {
	LogManager *log.LogManager
}

func Get() *HandlerFuncFactory {
	if handlerFuncFactory == nil {
		handlerFuncFactory = &HandlerFuncFactory{
			LogManager: log.Get(),
		}
	}

	return handlerFuncFactory
}

func (hff *HandlerFuncFactory) Produce(op operation.Operation) func(http.ResponseWriter, *http.Request) {
	switch op {
	case operation.Login:
		return http.HandlerFunc(api.Login)
	case operation.Save:
		return middleware.Middleware(http.HandlerFunc(api.Save))
	case operation.Retrieve:
		return middleware.Middleware(http.HandlerFunc(api.Retrieve))
	case operation.Update:
		return middleware.Middleware(http.HandlerFunc(api.Update))
	default:
		//todo log
		//log.Println(invalidOperationMessage)
		return nil
	}
}
