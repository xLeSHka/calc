package grpc

import (
	"context"
	"net/http"

	"github.com/xLeSHka/calculator/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Интерсептр для логирования запросов на grpc
func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any,err error) {
		l.Info(ctx,"request started",zap.String("method",info.FullMethod))
		return handler(ctx,req)
	}
}
// Мидлваря для логирования запросов на gateway
func LoggingMiddleware(handler http.Handler, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "request started", zap.String("method", r.URL.Path))
		handler.ServeHTTP(w, r)
	})
}
