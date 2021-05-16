package logger

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l := NewLogger()
		next.ServeHTTP(w, req.WithContext(NewContext(req.Context(), l)))
	})
}
