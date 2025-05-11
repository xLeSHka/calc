package users

import (
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"github.com/xLeSHka/calc/internal/utils/validator"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Router struct {
	Router       *routers.Routers
	usersService service.UsersService
	Log          *zap.Logger
	cryptoKey    []byte
	validator    *validator.Validator
}
type FxOpts struct {
	fx.In
	Routers      *routers.Routers
	Logger       *zap.Logger
	UsersService service.UsersService
	Config       config.Config
	Validator    *validator.Validator
}

func ClientRoute(opts FxOpts) *Router {
	router := &Router{
		Router:       opts.Routers,
		Log:          opts.Logger,
		usersService: opts.UsersService,
		cryptoKey:    []byte(opts.Config.CryptoKey),
		validator:    opts.Validator,
	}
	opts.Routers.ClientRoute.POST("/calculate", router.CreateExpression)
	opts.Routers.ClientRoute.GET("/expressions", router.GetExpressions)
	opts.Routers.ClientRoute.GET("/expressions/:id", router.GetExpression)
	//routers.Public.GET("/internal/task", router.GetTask)
	//routers.Public.POST("/internal/task", router.PostResult)
	return router
}
