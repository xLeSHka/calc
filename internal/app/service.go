package app

import (
	agentService "github.com/xLeSHka/calc/internal/orchestrator/service/agent"
	usersService "github.com/xLeSHka/calc/internal/orchestrator/service/users"
	"go.uber.org/fx"
)

var Services = fx.Provide(
	usersService.New,
	agentService.New,
)
