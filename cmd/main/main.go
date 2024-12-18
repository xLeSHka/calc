package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/xLeSHka/calc/internal/config"
	"github.com/xLeSHka/calc/internal/server"
	"github.com/xLeSHka/calc/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	//инициализируем логер и пихаем его в контекст
	ctx := context.Background()
	mainLogger := logger.New()
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	//загружаем конфиг
	cfg, err := config.New()
	if err != nil {
		mainLogger.Error(ctx, "error load config", zap.String("error message", err.Error()))
		return
	}

	//создаем экземпляр приложения
	server, err := server.New(ctx, cfg.RestServerPort)
	if err != nil {
		mainLogger.Error(ctx, "failed create new grpc server", zap.String("error message", err.Error()))
		return
	}
	//
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//задаем порт для сервера
	port := fmt.Sprintf(":%d", cfg.RestServerPort)

	//запускаем сервер
	go func() {
		if err := server.Start(port); err != nil && err != http.ErrServerClosed {
			mainLogger.Error(ctx, "shutting down the server", zap.String("Error:", err.Error()))
			return
		}
		mainLogger.Info(ctx, "Server started on port", zap.Int("port", cfg.RestServerPort))
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Плавно завершаем работу сервера
	if err := server.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}
	mainLogger.Info(ctx, "Server stopped")
}
