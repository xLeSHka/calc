package agent

import (
	"github.com/xLeSHka/calc/internal/agent"
	config2 "github.com/xLeSHka/calc/internal/pkg/config"
	logger2 "github.com/xLeSHka/calc/internal/pkg/logger"
	"go.uber.org/fx"
)

var Agent = fx.Option(
	fx.Provide(
		logger2.New,
		config2.New,
		agent.New,
	),
)
