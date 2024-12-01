package logger

import (
	"context"

	"go.uber.org/zap"
)

var (
	LoggerKey = "logger"
	RequestID = "requestID"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}
type logger struct {
	logger *zap.Logger
}

// Запись лога на уровне Info
func (l logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	//Записываем RequestID для отслеживания запроса
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.logger.Info(msg, fields...)
}

// Запись лога на уровне Error
func (l logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	//Записываем RequestID для отслеживания запроса
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.logger.Error(msg, fields...)
}

// Функция конструктор для логгера
func New() Logger {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	return &logger{
		logger: zapLogger,
	}
}

// Функция получения логера из контекста
func GetLoggerFromCtx(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)
}
