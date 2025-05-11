package app

import (
	repositoryAgent "github.com/xLeSHka/calc/internal/orchestrator/repository/agent"
	"github.com/xLeSHka/calc/internal/orchestrator/repository/redisRepo"
	repositoryUsers "github.com/xLeSHka/calc/internal/orchestrator/repository/users"
	"go.uber.org/fx"
)

var Repositories = fx.Provide(
	repositoryUsers.New,
	repositoryAgent.New,
	redisRepo.New,
)
