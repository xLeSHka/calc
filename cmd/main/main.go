package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xLeSHka/calculator/interenal/application"
	"github.com/xLeSHka/calculator/interenal/config"
	"github.com/xLeSHka/calculator/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	//инициализируем логер и пихаем его в контекст
	mainLogger := logger.New()
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	//загружаем конфиг
	cfg := config.New()
	if cfg == nil {
		log.Fatal("failed load config")
	}
	//создаем экземпляр приложения
	app := application.New(ctx, cfg.ServerPort)

	//создаем канал чтобы слушать сигналы системы interrupt и terminate
	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	//запускаем сервер
	go func() {
		if err := app.Start(ctx); err != nil {
			mainLogger.Info(ctx, "start server", zap.String("Error:", err.Error()))
		}
	}()
	<-graceCh

	//Плавно завершаем работу сервера
	if err := app.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}
