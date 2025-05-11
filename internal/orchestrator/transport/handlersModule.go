package transport

import (
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers/public"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers/users"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	"github.com/xLeSHka/calc/internal/utils/validator"
	"go.uber.org/fx"
)

var HttpModule = fx.Module("httpHandlers",
	fx.Provide(
		routers.CreateRouter,
		validator.New,
		fx.Private,
	),
	fx.Invoke(
		public.PublicRoute,
		users.ClientRoute,
	),
)
