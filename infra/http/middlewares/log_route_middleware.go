package middlewares

import (
	"net/http"
	"time"

	"github.com/diegodesousas/clean-boilerplate-go/infra/logger"
)

func LogRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(t time.Time) {
			logger.FromContext(req.Context()).
				Infof(
					"http: method=%s path=%s time=%dms",
					req.Method,
					req.URL.Path,
					time.Since(t).Milliseconds(),
				)
		}(time.Now())

		next.ServeHTTP(w, req)
	})
}
