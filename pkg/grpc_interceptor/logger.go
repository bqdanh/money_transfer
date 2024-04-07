package grpc_interceptor

import (
	"context"

	"github.com/bqdanh/money_transfer/pkg/logger"
	"github.com/google/uuid"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ZapLoggerUnaryServerInterceptor(shouldLog func(fullMethodName string, err error) bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := zap.S().With("trace_id", uuid.NewString())
		ctx = logger.ContextWithLogger(ctx, l)

		return grpc_zap.UnaryServerInterceptor(l.Desugar(), grpc_zap.WithDecider(shouldLog))(ctx, req, info, handler)
	}
}

func DefaultShouldLog(methodsNoLog map[string]struct{}) func(fullMethodName string, err error) bool {
	return func(fullMethodName string, err error) bool {
		if err != nil {
			return true
		}
		if _, ok := methodsNoLog[fullMethodName]; ok {
			return false
		}

		return true
	}
}
