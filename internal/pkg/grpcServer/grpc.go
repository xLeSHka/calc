package grpcServer

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	router "github.com/xLeSHka/calc/internal/orchestrator/transport/grpc"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

type GRPCServer struct {
	Server   *grpc.Server
	Listener net.Listener
}
type FxOpts struct {
	fx.In
	Config       config.Config
	Log          *zap.Logger
	AgentService service.AgentService
}

func New(opts FxOpts, lc fx.Lifecycle) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.Config.ServerHost, opts.Config.ServerPort+1))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterCalcServiceServer(grpcServer, router.New(router.FxOpts{AgentService: opts.AgentService, Logger: opts.Log}))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := grpcServer.Serve(lis)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("Server shutdown", zap.Error(err))
				}
			}()
			log.Info("Server started at address", zap.String("address", fmt.Sprintf("%s:%d", opts.Config.ServerHost, opts.Config.ServerPort+1)))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			grpcServer.GracefulStop()
			log.Error("Server stopped")
			return nil
		},
	})
	return &GRPCServer{Server: grpcServer, Listener: lis}, nil
}
