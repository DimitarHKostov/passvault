package handler_func_factory

import (
	"net/http"
	"passvault/pkg/api"
	"passvault/pkg/operation"
)

type HandlerFuncFactory struct {
}

var (
	handlerFuncFactory *HandlerFuncFactory
)

func Get() *HandlerFuncFactory {
	if handlerFuncFactory == nil {
		handlerFuncFactory = &HandlerFuncFactory{}
	}

	return handlerFuncFactory
}

func (hff *HandlerFuncFactory) Produce(op operation.Operation) func(http.ResponseWriter, *http.Request) {
	switch op {
	case operation.Login:
		return http.HandlerFunc(api.Login)
	case operation.Save:
		return api.Middleware(http.HandlerFunc(api.Save))
	case operation.Retrieve:
		return api.Middleware(http.HandlerFunc(api.Retrieve))
	case operation.Update:
		return api.Middleware(http.HandlerFunc(api.Update))
	default:
		return nil
	}
}
