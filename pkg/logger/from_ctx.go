package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextLoggerKey struct{}

func ContextWithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, contextLoggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	vl := ctx.Value(contextLoggerKey{})
	if vl == nil {
		return zap.S()
	}
	logger, ok := vl.(*zap.SugaredLogger)
	if !ok {
		return zap.S()
	}
	return logger
}
