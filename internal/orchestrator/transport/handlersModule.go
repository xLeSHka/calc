package transport

import (
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	"go.uber.org/fx"
)

var HttpModule = fx.Module("httpHandlers",
	fx.Provide(
		routers.CreateRouter,
	),
	fx.Invoke(
		handlers.SetUpRouter,
	),
)
