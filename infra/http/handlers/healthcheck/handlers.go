package healthcheck

import (
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
)

func Readiness(conn database.Conn) server.ActionHandler {
	return func(w http.ResponseWriter, req *http.Request) error {
		ctx := req.Context()

		var dest string
		if err := conn.GetContext(ctx, &dest, "SELECT 1"); err != nil {
			return err
		}

		_, err := w.Write([]byte("Ok"))

		return err
	}
}

func Liveness(w http.ResponseWriter, req *http.Request) error {
	_, err := w.Write([]byte("Ok"))

	return err
}
