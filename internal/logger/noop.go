package logger

import "context"

type NopLogger struct{}

func NewNoopLogger() *NopLogger {
	return &NopLogger{}
}

func (NopLogger) Debugf(_ context.Context, _ string, _ ...any) {}
func (NopLogger) Infof(_ context.Context, _ string, _ ...any)  {}
