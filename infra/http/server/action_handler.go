package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/domain/errors"
	"github.com/diegodesousas/clean-boilerplate-go/infra/logger"
)

type ActionHandler func(w http.ResponseWriter, req *http.Request) error

func ActionErrorHandler(handler ActionHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)

		if err != nil {
			handleError(req.Context(), w, err)
		}
	})
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	entry := logger.
		FromContext(ctx).
		WithField("app.type", "action-error-handler")

	switch err.(type) {
	case *errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
		return

	case *errors.ValidationErrors:
		bytes, _ := json.Marshal(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(bytes)
		return

	case *errors.Conflict:
		bytes, _ := json.Marshal(err)
		w.WriteHeader(http.StatusConflict)
		w.Write(bytes)
		return

	case *errors.Unprocessable:
		bytes, _ := json.Marshal(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(bytes)
		return

	default:
		entry.
			WithError(err).
			Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
