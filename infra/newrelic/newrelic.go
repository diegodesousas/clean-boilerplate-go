package newrelic

import (
	"net/http"

	"github.com/diegodesousas/clean-boilerplate-go/infra/http/server"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
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

func NewMonitorWrapper() server.MonitorWrapper {
	nrApp, err := New()
	if err != nil {
		return nil
	}

	return func(path string, handler http.Handler) http.Handler {
		_, h := newrelic.WrapHandle(nrApp, path, handler)
		return h
	}
}
