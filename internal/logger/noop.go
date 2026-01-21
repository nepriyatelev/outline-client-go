package logger

import "context"

type nopLogger struct{}

func NewNoopLogger() *nopLogger {
	return &nopLogger{}
}

func (nopLogger) Debugf(_ context.Context, _ string, _ ...any) {}
func (nopLogger) Infof(_ context.Context, _ string, _ ...any)  {}
