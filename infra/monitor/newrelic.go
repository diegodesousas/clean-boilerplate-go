package monitor

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
)

type NewRelicWrapper func(path string, handler http.Handler) http.Handler

var (
	NewRelicWrapperDefault = func(path string, handler http.Handler) http.Handler {
		return handler
	}
)

func New() (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(viper.GetString("NEWRELIC_NAME")),
		newrelic.ConfigLicense(viper.GetString("NEWRELIC_LICENSE_KEY")),
		newrelic.ConfigEnabled(viper.GetBool("NEWRELIC_ENABLE")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func NewNewRelicWrapper(app *newrelic.Application) NewRelicWrapper {
	return func(path string, handler http.Handler) http.Handler {
		_, h := newrelic.WrapHandle(app, path, handler)
		return h
	}
}
