package transport

import (
	"github.com/xLeSHka/calc/internal/app/transport/handlers"
	"github.com/xLeSHka/calc/internal/app/transport/routers"
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
