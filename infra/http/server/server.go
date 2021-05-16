package server

import (
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/infra/http/middlewares"
	"github.com/diegodesousas/clean-boilerplate-go/infra/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

type Server struct {
	http.Server
	router *httprouter.Router
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
		router: httprouter.New(),
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

	return server
}

func (s *Server) ListenAndServe() error {
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Route(r Route) {
	s.router.Handler(r.Method, r.Path, middlewares.NewRelicWrapper(r.Path, r.Handler))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, req)
}
