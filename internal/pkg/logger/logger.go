package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var ErrFailedBuildLogger = errors.New("logger: failed build logger")

func New() (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%w, error: %w", ErrFailedBuildLogger, err)
	}
	defer zapLogger.Sync()
	return zapLogger, nil
}
