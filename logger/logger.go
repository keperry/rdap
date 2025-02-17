package logger

import (
	"context"
)

type Logger interface {
	Logf(ctx context.Context, format string, args ...interface{})
}

type NoopLogger struct{}

func (NoopLogger) Logf(ctx context.Context, format string, args ...interface{}) {}

type LoggerFunc func(ctx context.Context, format string, args ...interface{})

func (lf LoggerFunc) Logf(ctx context.Context, format string, args ...interface{}) {
	lf(ctx, format, args...)
}
