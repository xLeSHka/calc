package grpc

import (
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	"go.uber.org/zap"
)

type Router struct {
	proto.UnimplementedCalcServiceServer
	agentService service.AgentService
	log          *zap.Logger
}
type FxOpts struct {
	AgentService service.AgentService
	Logger       *zap.Logger
}

func New(opts FxOpts) *Router {
	return &Router{agentService: opts.AgentService, log: opts.Logger}
}
