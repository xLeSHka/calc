package orchestrator

import (
	"github.com/xLeSHka/calc/internal/app"
	"github.com/xLeSHka/calc/internal/orchestrator/transport"
	"github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/grpcServer"
	"github.com/xLeSHka/calc/internal/pkg/http"
	"github.com/xLeSHka/calc/internal/pkg/logger"
	"github.com/xLeSHka/calc/internal/pkg/postgres"
	"github.com/xLeSHka/calc/internal/pkg/rdb"
	"github.com/xLeSHka/calc/internal/utils/jwt"
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
		rdb.New,
		jwt.New,
		//fx.Annotate(agent.New,
		//	fx.As(new(repository.UsersRepository)),
		//),
		//calculator.New,
		//fx.Annotate(service.New,
		//	fx.As(new(routers.Service)),
		//),
	),
	app.Repositories,
	fx.Provide(
		calculator.New,
	),
	app.Services,
	fx.Invoke(
		postgres.MigrateDB,
	),
	transport.HttpModule,
	fx.Invoke(
		grpcServer.New,
	),
)
