package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegodesousas/clean-boilerplate-go/infra/http/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

var DefaultMonitorWrapper = func(path string, handler http.Handler) http.Handler {
	return handler
}

type (
	Config         func(server *Server)
	MonitorWrapper func(path string, handler http.Handler) http.Handler
)

type Server struct {
	http.Server
	routes         []Route
	router         *httprouter.Router
	monitorWrapper MonitorWrapper
}

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
		router:         httprouter.New(),
		monitorWrapper: DefaultMonitorWrapper,
	}

	server.Server.Handler = middlewares.Middlewares(
		server,
		middlewares.PanicRecoveryMiddleware,
		middlewares.Logger,
		middlewares.LogRouteMiddleware,
	)

	for _, config := range configs {
		config(server)
	}

	server.buildRoutes()

	return server
}

type ShutdownHandler func(ctx context.Context) error

func (s *Server) ListenAndServe(ctx context.Context, handlers ...ShutdownHandler) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print(err)
		}

		interrupt <- syscall.SIGTERM
	}()

	<-interrupt
	log.Println("shutdown application")
	log.Println("shutdown http server")
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("http: %s", err)
	}

	for _, handler := range handlers {
		err := handler(ctx)
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

func (s *Server) Route(r Route) {
	s.routes = append(s.routes, r)
}

func (s *Server) buildRoutes() {
	for _, r := range s.routes {
		s.router.Handler(r.Method, r.Path, s.monitorWrapper(r.Path, r.Handler))
	}

	s.routes = []Route{}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, req)
}

func WithMonitorWrapper(wrapper MonitorWrapper) Config {
	return func(server *Server) {
		if wrapper == nil {
			return
		}

		server.monitorWrapper = wrapper
	}
}
