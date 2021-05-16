package server

import (
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/infra/monitor"

	"github.com/diegodesousas/clean-boilerplate-go/infra/http/middlewares"
	"github.com/diegodesousas/clean-boilerplate-go/infra/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
)

type Server struct {
	http.Server
	routes    []Route
	router    *httprouter.Router
	nrWrapper monitor.NewRelicWrapper
}

type Config func(server *Server)

type Route struct {
	Path    string
	Method  string
	Handler http.Handler
}

func NewServer(configs ...Config) *Server {
	server := &Server{
		Server: http.Server{
			Addr: ":" + viper.GetString("HTTP_PORT"),
		},
		router:    httprouter.New(),
		nrWrapper: monitor.NewRelicWrapperDefault,
	}

	server.Server.Handler = middlewares.Middlewares(
		server,
		middlewares.PanicRecoveryMiddleware,
		logger.Middleware,
		middlewares.LogRouteMiddleware,
	)

	for _, config := range configs {
		config(server)
	}

	server.buildRoutes()

	return server
}

func (s *Server) ListenAndServe() error {
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Route(r Route) {
	s.routes = append(s.routes, r)
}

func (s *Server) buildRoutes() {
	for _, r := range s.routes {
		s.router.Handler(r.Method, r.Path, s.nrWrapper(r.Path, r.Handler))
	}

	s.routes = []Route{}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, req)
}

func WithNewRelicWrapper(app *newrelic.Application) Config {
	return func(server *Server) {
		server.nrWrapper = monitor.NewNewRelicWrapper(app)
	}
}
