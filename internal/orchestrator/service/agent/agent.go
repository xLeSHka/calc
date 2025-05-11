package agentService

import (
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AgentService struct {
	Log             *zap.Logger
	AgentRepository repository.AgentRepository
	Calculator      *calculator.Calculator
}
type FxOpts struct {
	fx.In
	Log             *zap.Logger
	AgentRepository repository.AgentRepository
	Calculator      *calculator.Calculator
}

func New(opts FxOpts) service.AgentService {
	return &AgentService{
		Log:             opts.Log,
		AgentRepository: opts.AgentRepository,
		Calculator:      opts.Calculator,
	}
}
