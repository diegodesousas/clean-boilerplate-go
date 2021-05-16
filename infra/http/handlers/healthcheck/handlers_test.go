package healthcheck

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
	mockDatabase "github.com/diegodesousas/clean-boilerplate-go/test/mock/infra/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadiness(t *testing.T) {
	type args struct {
		conn       database.Conn
		w          http.ResponseWriter
		req        *http.Request
		body       string
		statusCode int
	}

	type build func(*testing.T) args

	tests := []struct {
		name        string
		prepareArgs build
		want        error
	}{
		{
			name: "readiness success",
			prepareArgs: func(t *testing.T) args {
				ctrl := gomock.NewController(t)
				rr := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, "/readiness", nil)

				mockConn := mockDatabase.NewMockConn(ctrl)
				mockConn.
					EXPECT().
					GetContext(gomock.Any(), gomock.Any(), "SELECT 1").
					Return(nil)

				return args{
					conn:       mockConn,
					w:          rr,
					req:        req,
					body:       "Ok",
					statusCode: http.StatusOK,
				}
			},
			want: nil,
		},
		{
			name: "readiness fail",
			prepareArgs: func(t *testing.T) args {
				ctrl := gomock.NewController(t)
				rr := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, "/readiness", nil)

				mockConn := mockDatabase.NewMockConn(ctrl)
				mockConn.
					EXPECT().
					GetContext(gomock.Any(), gomock.Any(), "SELECT 1").
					Return(errors.New("database connection error"))

				return args{
					conn:       mockConn,
					w:          rr,
					req:        req,
					body:       "",
					statusCode: http.StatusInternalServerError,
				}
			},
			want: errors.New("database connection error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := assert.New(t)
			args := tt.prepareArgs(t)

			handler := Readiness(args.conn)

			server.ActionErrorHandler(handler).ServeHTTP(args.w, args.req)

			rr := args.w.(*httptest.ResponseRecorder)

			expect.Equal(args.statusCode, rr.Result().StatusCode, "status code is not ok")

			bb, err := ioutil.ReadAll(rr.Result().Body)
			expect.Nil(err, "read body error is not nil")
			expect.Equal(args.body, string(bb), "string body not match with expected")
		})
	}
}

func TestLiveness(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		req        *http.Request
		body       string
		statusCode int
	}

	type build func(*testing.T) args

	tests := []struct {
		name        string
		prepareArgs build
		want        error
	}{
		{
			name: "liveness success",
			prepareArgs: func(t *testing.T) args {
				rr := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, "/liveness", nil)

				return args{
					w:          rr,
					req:        req,
					body:       "Ok",
					statusCode: http.StatusOK,
				}
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := assert.New(t)
			args := tt.prepareArgs(t)

			server.ActionErrorHandler(Liveness).ServeHTTP(args.w, args.req)

			rr := args.w.(*httptest.ResponseRecorder)

			expect.Equal(args.statusCode, rr.Result().StatusCode, "status code is not ok")

			bb, err := ioutil.ReadAll(rr.Result().Body)
			expect.Nil(err, "read body error is not nil")
			expect.Equal(args.body, string(bb), "string body not match with expected")
		})
	}
}
