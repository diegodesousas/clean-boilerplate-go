package logger

import (
	"context"

	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/sirupsen/logrus"
)

type key string

var loggerKey key = "logger"

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(nrlogrusplugin.ContextFormatter{})

	return &Logger{
		l,
	}
}

func NewContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(loggerKey).(*Logger)

	if !ok {
		logger = NewLogger()
	}

	return logger.WithContext(ctx)
}
