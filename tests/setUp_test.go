package tests

import (
	"context"
	"errors"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pressly/goose"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/app/orchestrator"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	repositoryAgent "github.com/xLeSHka/calc/internal/orchestrator/repository/agent"
	"github.com/xLeSHka/calc/internal/orchestrator/repository/redisRepo"
	repositoryUsers "github.com/xLeSHka/calc/internal/orchestrator/repository/users"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	agentService "github.com/xLeSHka/calc/internal/orchestrator/service/agent"
	usersService "github.com/xLeSHka/calc/internal/orchestrator/service/users"
	router "github.com/xLeSHka/calc/internal/orchestrator/transport/grpc"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers/public"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers/users"
	routers2 "github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	cache2 "github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	config2 "github.com/xLeSHka/calc/internal/pkg/config"
	counter2 "github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	grpcServer2 "github.com/xLeSHka/calc/internal/pkg/grpcServer"
	http3 "github.com/xLeSHka/calc/internal/pkg/http"
	logger2 "github.com/xLeSHka/calc/internal/pkg/logger"
	"github.com/xLeSHka/calc/internal/pkg/postgres"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"github.com/xLeSHka/calc/internal/utils/validator"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(orchestrator.Orchestrator)
	assert.Nil(t, err)
}

var quit = make(chan os.Signal, 1)
var config config2.Config
var logger *zap.Logger
var db *gorm.DB
var http2 *gin.Engine
var cache *cache2.Cache
var counter *counter2.Counter
var rdb *redis.Client
var RedisRepository repository.RedisRepository
var AgentRepositpry repository.AgentRepository
var UsersRepository repository.UsersRepository
var JWT *jwt.JWT
var Validator *validator.Validator
var UsersService service.UsersService
var AgentService service.AgentService
var Calculator *calculator.Calculator
var routers *routers2.Routers
var GRPCServer *grpcServer2.GRPCServer
var profile1JWT string
var unknownJWT string
var profile2JWT string
var profile3JWT string
var client proto.CalcServiceClient

func init() {
	gin.SetMode(gin.TestMode)
	var err error
	logger, err = logger2.New()
	if err != nil {
		log.Fatal(err)
	}
	config, err = config2.New()
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	db, err = postgres.New(config)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisHost + ":" + strconv.Itoa(int(config.RedisPort)),
		DB:   0,
	})

	err = rdb.Ping(context.Background()).Err()

	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	counter = counter2.New()
	http2 = gin.Default()
	http2.Use(gin.Recovery())
	http2.Use(http3.CORSMiddleware())
	Validator = validator.New()
	JWT = jwt.New(config)
	err = postgres.MigrateDB(db)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	cache = cache2.New(config)
	RedisRepository = redisRepo.New(rdb)
	UsersRepository = repositoryUsers.New(db)
	AgentRepositpry = repositoryAgent.New(db)
	Calculator = calculator.New(calculator.FxOpts{
		Log:     logger,
		Repo:    AgentRepositpry,
		Times:   cache,
		Counter: counter,
	})
	UsersService = usersService.New(usersService.FxOpts{
		UsersRepository: UsersRepository,
		Log:             logger,
		Calculator:      Calculator,
		JWT:             JWT,
		RDB:             RedisRepository,
		Config:          config,
	})
	AgentService = agentService.New(agentService.FxOpts{
		Log:             logger,
		AgentRepository: AgentRepositpry,
		Calculator:      Calculator,
	})

	routers = routers2.CreateRouter(http2, JWT, rdb)
	users.ClientRoute(users.FxOpts{
		Routers:      routers,
		Logger:       logger,
		UsersService: UsersService,
		Config:       config,
		Validator:    Validator,
	})
	public.PublicRoute(public.FxOpts{
		Routers:      routers,
		Logger:       logger,
		UsersService: UsersService,
		Config:       config,
		Validator:    Validator,
	})

}
func setUp() (func(), chan os.Signal, error) {
	var err error
	profile1JWT, err = JWT.CreateToken(jwt2.MapClaims{
		"id": profile1.ID,
	}, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatal(err)
	}
	unknownJWT, err = JWT.CreateToken(jwt2.MapClaims{
		"id": uuid.New(),
	}, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatal(err)
	}
	profile2JWT, err = JWT.CreateToken(jwt2.MapClaims{
		"id": profile2.ID,
	}, time.Now().Add(time.Hour*24*7))
	profile3JWT, err = JWT.CreateToken(jwt2.MapClaims{
		"id": uuid.New(),
	}, time.Now().Add(time.Hour*24*7))
	RedisRepository.SaveToken(context.TODO(), profile1.ID, profile1JWT)
	RedisRepository.SaveToken(context.TODO(), profile2.ID, profile2JWT)
	logger.Info("starting server")
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, nil, err
	}
	db3, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	if err := goose.Up(db3, "migrations"); err != nil {
		return nil, nil, err
	}
	db.Exec("DELETE FROM public.expressions;")
	db.Exec("DELETE FROM public.users;")

	counter.Restart()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort+3),
		Handler: http2,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen:", zap.Error(err))
		}
	}()
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort+4))
		if err != nil {
			logger.Fatal("", zap.Error(err))
		}
		s := grpc.NewServer()
		proto.RegisterCalcServiceServer(s, router.New(router.FxOpts{AgentService: AgentService, Logger: logger}))

		GRPCServer = &grpcServer2.GRPCServer{
			Server:   s,
			Listener: lis,
		}
		err = GRPCServer.Server.Serve(GRPCServer.Listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Fatal("server shutdown", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	fn := func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		GRPCServer.Server.GracefulStop()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("server shutdown", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			logger.Info("time out 1s")
		}
		logger.Info("server exiting")
	}
	logger.Info("server started", zap.String("host", srv.Addr))
	return fn, quit, nil
}

type SendExpression struct {
	Expression string `json:"expression"`
}

var profile1 models.User = models.User{
	ID:    uuid.New(),
	Login: "user_1",
}
var profile2 models.User = models.User{
	ID:    uuid.New(),
	Login: "user_2",
}

var (
	r                      = 4.0
	r2                     = 5.43311
	Expr models.Expression = models.Expression{
		ID:         1111,
		UserID:     profile1.ID,
		Expression: "2+2",
		Status:     "Solved",
		Result:     &r,
	}
	Expr2 models.Expression = models.Expression{
		ID:         1112,
		UserID:     profile1.ID,
		Expression: "log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)",
		Status:     "Solved",
		Result:     &r2,
	}
)

func solveExpr() {
	var cErr *customError.CustomError = nil
	var j *models.Task = nil
	for {
		j, cErr = AgentService.GetTask()
		if cErr != nil {
			break
		}
		logger.Info("", zap.Int64("id", j.ID), zap.Int64("expression id", j.ExpressionID))
		switch j.Operation {
		case models.Addition:
			AgentService.SetResult(j.ID, j.ExpressionID, j.Arg1+j.Arg2, nil)
		case models.Subtraction:
			AgentService.SetResult(j.ID, j.ExpressionID, j.Arg1-j.Arg2, nil)
		case models.Multiplication:
			AgentService.SetResult(j.ID, j.ExpressionID, j.Arg1*j.Arg2, nil)
		case models.Division:
			if j.Arg2 == 0 {
				err := "Division by zero"
				AgentService.SetResult(j.ID, j.ExpressionID, 0.0, &err)
			} else {
				AgentService.SetResult(j.ID, j.ExpressionID, j.Arg1/j.Arg2, nil)

			}
		case models.Exponentiation:
			AgentService.SetResult(j.ID, j.ExpressionID, math.Pow(j.Arg1, j.Arg2), nil)
		case models.UnaryMinus:
			AgentService.SetResult(j.ID, j.ExpressionID, -j.Arg1, nil)
		case models.Logarithm:
			if j.Arg1 <= 0 || j.Arg1 == 1 {
				errMsg := "log not defined"
				AgentService.SetResult(j.ID, j.ExpressionID, 0.0, &errMsg)
			} else if j.Arg2 <= 0.0 {
				errMsg := "log out of domain"
				AgentService.SetResult(j.ID, j.ExpressionID, 0.0, &errMsg)
			} else {
				AgentService.SetResult(j.ID, j.ExpressionID, math.Log(j.Arg2)/math.Log(j.Arg1), nil)
			}
		case models.SquareRoot:
			if j.Arg1 < 0 {
				errMsg := "negative square"
				AgentService.SetResult(j.ID, j.ExpressionID, 0.0, &errMsg)
			} else {
				AgentService.SetResult(j.ID, j.ExpressionID, math.Sqrt(j.Arg1), nil)
			}
		}
	}
}

type PostTask struct {
	ID           int64   `json:"id"`
	ExpressionID int64   `json:"expression_id"`
	Result       float64 `json:"result"`
}
