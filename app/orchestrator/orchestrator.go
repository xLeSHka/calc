package orchestrator

import (
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/orchestrator/repository/db"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	"github.com/xLeSHka/calc/internal/orchestrator/transport"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	"github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/http"
	"github.com/xLeSHka/calc/internal/pkg/logger"
	"github.com/xLeSHka/calc/internal/pkg/postgres"
	"go.uber.org/fx"
)

var Orchestrator = fx.Options(
	fx.Provide(
		config.New,
		logger.New,
		http.New,
		counter.New,
		cache.New,
		postgres.New,
		fx.Annotate(db.New,
			fx.As(new(repository.Repo)),
		),
		calculator.New,
		fx.Annotate(service.New,
			fx.As(new(routers.Service)),
		),
	),
	fx.Invoke(
		postgres.MigrateDB,
	),
	transport.HttpModule,
)
