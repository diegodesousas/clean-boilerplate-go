package middlewares

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
)

var defaultApp *newrelic.Application

func InitNewRelic() error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(viper.GetString("NEWRELIC_NAME")),
		newrelic.ConfigLicense(viper.GetString("NEWRELIC_LICENSE_KEY")),
		newrelic.ConfigEnabled(viper.GetBool("NEWRELIC_ENABLE")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return err
	}

	defaultApp = app

	return nil
}

func NewRelicWrapper(path string, handler http.Handler) http.Handler {
	_, h := newrelic.WrapHandle(defaultApp, path, handler)
	return h
}
