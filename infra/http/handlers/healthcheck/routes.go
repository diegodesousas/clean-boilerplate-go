package healthcheck

import (
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
)

func Routes(conn database.Conn) server.Config {
	return func(s *server.Server) {
		s.Route(server.Route{
			Method:  http.MethodGet,
			Path:    "/readiness",
			Handler: server.ActionErrorHandler(Readiness(conn)),
		})

		s.Route(server.Route{
			Method:  http.MethodGet,
			Path:    "/liveness",
			Handler: server.ActionErrorHandler(Liveness),
		})
	}
}
