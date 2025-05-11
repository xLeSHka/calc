package usersService

import (
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UsersService struct {
	Log             *zap.Logger
	UsersRepository repository.UsersRepository
	Calculator      *calculator.Calculator
	JWT             *jwt.JWT
	RDB             repository.RedisRepository
	cryptoKey       []byte
}
type FxOpts struct {
	fx.In
	UsersRepository repository.UsersRepository
	Log             *zap.Logger
	Calculator      *calculator.Calculator
	JWT             *jwt.JWT
	RDB             repository.RedisRepository
	Config          config.Config
}

func New(
	fxOpts FxOpts,
) service.UsersService {
	return &UsersService{
		Log:             fxOpts.Log,
		UsersRepository: fxOpts.UsersRepository,
		Calculator:      fxOpts.Calculator,
		JWT:             fxOpts.JWT,
		RDB:             fxOpts.RDB,
		cryptoKey:       []byte(fxOpts.Config.CryptoKey),
	}
}
