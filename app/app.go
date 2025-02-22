package app

import (
	"github.com/xLeSHka/calc/internal/app/repository"
	"github.com/xLeSHka/calc/internal/app/repository/mock"
	"github.com/xLeSHka/calc/internal/app/service"
	"github.com/xLeSHka/calc/internal/app/transport"
	"github.com/xLeSHka/calc/internal/app/transport/routers"
	"github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/http"
	"github.com/xLeSHka/calc/internal/pkg/logger"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		config.New,
		logger.New,
		http.New,
		counter.New,
		cache.New,
		fx.Annotate(mock.New,
			fx.As(new(repository.Repo)),
		),
		calculator.New,
		fx.Annotate(service.New,
			fx.As(new(routers.Service)),
		),
	),
	transport.HttpModule,
)
