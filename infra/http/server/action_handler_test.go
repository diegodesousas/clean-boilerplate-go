package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	domainerrors "github.com/diegodesousas/clean-boilerplate-go/domain/errors"
)

func TestActionErrorHandler(t *testing.T) {
	type args struct {
		handler ActionHandler
		w       *httptest.ResponseRecorder
		req     *http.Request
	}
	type expected struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "should run action handle function executed successfully",
			args: args{
				handler: func(w http.ResponseWriter, req *http.Request) error {
					return nil
				},
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
			expected: expected{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name: "should return internal server error",
			args: args{
				handler: func(w http.ResponseWriter, req *http.Request) error {
					return errors.New("error test")
				},
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				body:       "",
			},
		},
		{
			name: "should return not found status",
			args: args{
				handler: func(w http.ResponseWriter, req *http.Request) error {
					return domainerrors.NewEntityNotFound("test")
				},
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				body:       "",
			},
		},
		{
			name: "should return unprocessable entity status",
			args: args{
				handler: func(w http.ResponseWriter, req *http.Request) error {
					err := domainerrors.NewValidationErrors()
					err.Add("test", domainerrors.Required)
					return err
				},
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				body:       "{\"errors\":{\"test\":{\"messages\":{\"required\":\"The field test is required\"}}}}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ActionErrorHandler(tt.args.handler).ServeHTTP(tt.args.w, tt.args.req)

			res := tt.args.w.Result()
			if res.StatusCode != tt.expected.statusCode {
				t.Errorf("expected status code %d; got %d;", tt.expected.statusCode, tt.args.w.Code)
			}

			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(tt.expected.body, string(bodyBytes)) {
				t.Errorf("expected body %s; got %s;", tt.expected.body, string(bodyBytes))
				return
			}
		})
	}
}
