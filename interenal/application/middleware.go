package application

import (
	"context"
	"net/http"

	"github.com/xLeSHka/calculator/pkg/logger"
	"go.uber.org/zap"
)

// Мидлваря ля логирования запросов
func LoggingMiddleware(handler http.Handler, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "request started", zap.String("method", r.URL.Path))
		handler.ServeHTTP(w, r)
	})
}
