package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKeyLogger struct{}

func init() {
	formatter := logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(&formatter)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
}

func Ctx(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(ctxKeyLogger{}).(*logrus.Entry)
	if !ok {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return logger
}

func WithField(ctx context.Context, key string, value interface{}) (context.Context, *logrus.Entry) {
	logger := Ctx(ctx).WithField(key, value)
	return context.WithValue(ctx, ctxKeyLogger{}, logger), logger
}

func WithFields(ctx context.Context, fields map[string]interface{}) (context.Context, *logrus.Entry) {
	logger := Ctx(ctx).WithFields(fields)
	ctx = context.WithValue(ctx, ctxKeyLogger{}, logger)
	return ctx, logger
}
