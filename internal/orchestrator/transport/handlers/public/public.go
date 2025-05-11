package public

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
	validator    *validator.Validator
	usersService service.UsersService
	Log          *zap.Logger
	cryptoKey    []byte
}
type FxOpts struct {
	fx.In
	Routers      *routers.Routers
	Logger       *zap.Logger
	UsersService service.UsersService
	Config       config.Config
	Validator    *validator.Validator
}

func PublicRoute(opts FxOpts) *Router {
	router := &Router{
		Router:       opts.Routers,
		Log:          opts.Logger,
		usersService: opts.UsersService,
		cryptoKey:    []byte(opts.Config.CryptoKey),
		validator:    opts.Validator,
	}
	opts.Routers.Public.POST("/register", router.Register)
	opts.Routers.Public.POST("/login", router.Login)
	return router
}
