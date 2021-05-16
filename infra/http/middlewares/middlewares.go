package middlewares

import (
	"net/http"
)

func Middlewares(main http.Handler, middlewares ...func(handler http.Handler) http.Handler) http.Handler {
	var h = main
	for i := range middlewares {
		h = middlewares[len(middlewares)-1-i](h)
	}

	return h
}
