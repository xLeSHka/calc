package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/xLeSHka/calculator/interenal/config"
	"github.com/xLeSHka/calculator/interenal/service"
	"github.com/xLeSHka/calculator/interenal/transport/grpc"
	"github.com/xLeSHka/calculator/pkg/logger"
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

	srv := service.New()

	//создаем экземпляр приложения
	grpcServer, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort, srv)
	if err != nil {
		mainLogger.Error(ctx, "failed create new grpc server", zap.String("error message", err.Error()))
		return
	}

	//создаем канал чтобы слушать сигналы системы interrupt и terminate
	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	//запускаем сервер
	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			mainLogger.Info(ctx, "start server", zap.String("Error:", err.Error()))
		}
	}()
	<-graceCh

	//Плавно завершаем работу сервера
	if err := grpcServer.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}
