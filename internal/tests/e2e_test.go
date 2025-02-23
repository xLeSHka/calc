package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/app/orchestrator"
	"github.com/xLeSHka/calc/internal/models"
	db2 "github.com/xLeSHka/calc/internal/orchestrator/repository/db"
	"github.com/xLeSHka/calc/internal/orchestrator/service"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/handlers"
	routers2 "github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	cache2 "github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	config2 "github.com/xLeSHka/calc/internal/pkg/config"
	counter2 "github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	http3 "github.com/xLeSHka/calc/internal/pkg/http"
	logger2 "github.com/xLeSHka/calc/internal/pkg/logger"
	"github.com/xLeSHka/calc/internal/pkg/postgres"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
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
var Repository *db2.Repository
var Calculator *calculator.Calculator
var Service *service.Service
var routers *routers2.Routers

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
	counter = counter2.New()
	http2 = gin.Default()
	http2.Use(gin.Recovery())
	http2.Use(http3.CORSMiddleware())
	err = postgres.MigrateDB(db)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	cache = cache2.New(config)
	Repository = db2.New(db)
	Calculator = calculator.New(logger, Repository, cache, counter)
	Service = service.New(Repository, logger, Calculator)

	routers = routers2.CreateRouter(http2)
	handlers.SetUpRouter(routers, logger, Service)

}
func setUp() (func(), chan os.Signal, error) {
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

	counter.Restart()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort+1),
		Handler: http2,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen:", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	fn := func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

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

func TestCreateExpression(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	defer solveExpr()
	type TestCreateExpression struct {
		Expression   SendExpression
		Expected     models.Expression
		name         string
		expectedCode int
	}
	tests := []TestCreateExpression{
		TestCreateExpression{
			name:       "Succes create expression",
			Expression: SendExpression{"2+2"},
			Expected: models.Expression{
				Expression: "2+2",
				Status:     "Solved",
				Result:     &r,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:       "Succes create expression 2",
			Expression: SendExpression{"log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)"},
			Expected: models.Expression{
				Expression: "log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)",
				Status:     "Solved",
				Result:     &r2,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:         "simple fail empty expression",
			Expression:   SendExpression{""},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:         "simple fail",
			Expression:   SendExpression{"1+1*"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:         "priority",
			Expression:   SendExpression{"2+2**2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:         "right paranthes",
			Expression:   SendExpression{"((2+2-*(2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:         "left paranthes",
			Expression:   SendExpression{"2+2)-2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:       "division by zero",
			Expression: SendExpression{"10/0"},
			Expected: models.Expression{
				Expression: "10/0",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:         "invaid operator",
			Expression:   SendExpression{"10&0"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			name:       "log bad req",
			Expression: SendExpression{"log(-2,8)"},
			Expected: models.Expression{
				Expression: "log(-2,8)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:       "log another bad req",
			Expression: SendExpression{"log(1,8)"},
			Expected: models.Expression{
				Expression: "log(1,8)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:       "log another bad req 2",
			Expression: SendExpression{"log(16,(-1))"},
			Expected: models.Expression{
				Expression: "log(16,(-1))",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:       "sqrt bad req",
			Expression: SendExpression{"sqrt(50-50-50)"},
			Expected: models.Expression{
				Expression: "sqrt(50-50-50)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			jsonData, _ := json.Marshal(test.Expression)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusCreated {
				for i := 0; i < 10; i++ {
					solveExpr()
					time.Sleep(100 * time.Millisecond)
				}
				var expr models.Expression
				result := db.Model(&models.Expression{}).Where("expression = ?", test.Expected.Expression).Order("id DESC").First(&expr)
				assert.Nil(t, result.Error)
				assert.Equal(t, test.Expected.Status, expr.Status)
				if test.Expected.Status == "Solved" {
					assert.Equal(t, *test.Expected.Result, *expr.Result)
				} else {
					assert.Nil(t, expr.Result)
				}
			}
		})
	}
}

var (
	r                      = 4.0
	r2                     = 5.43311
	Expr models.Expression = models.Expression{
		ID:         1111,
		Expression: "2+2",
		Status:     "Solved",
		Result:     &r,
	}
	Expr2 models.Expression = models.Expression{
		ID:         1112,
		Expression: "log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)",
		Status:     "Solved",
		Result:     &r2,
	}
)

func solveExpr() {
	var cErr *customError.CustomError = nil
	var j *models.Task = nil
	for {
		j, cErr = Service.GetTask()
		if cErr != nil {
			break
		}
		logger.Info("", zap.Int64("id", j.ID), zap.Int64("expression id", j.ExpressionID))
		switch j.Operation {
		case models.Addition:
			res := j.Arg1 + j.Arg2
			Service.SetResult(j.ID, j.ExpressionID, &res, nil)
		case models.Subtraction:
			res := j.Arg1 - j.Arg2
			Service.SetResult(j.ID, j.ExpressionID, &res, nil)
		case models.Multiplication:
			res := j.Arg1 * j.Arg2
			Service.SetResult(j.ID, j.ExpressionID, &res, nil)
		case models.Division:
			if j.Arg2 == 0 {
				err := "Division by zero"
				Service.SetResult(j.ID, j.ExpressionID, nil, &err)
			} else {
				res := j.Arg1 / j.Arg2
				Service.SetResult(j.ID, j.ExpressionID, &res, nil)

			}
		case models.Exponentiation:
			res := math.Pow(j.Arg1, j.Arg2)
			Service.SetResult(j.ID, j.ExpressionID, &res, nil)
		case models.UnaryMinus:
			res := -j.Arg1
			Service.SetResult(j.ID, j.ExpressionID, &res, nil)
		case models.Logarithm:
			if j.Arg1 <= 0 || j.Arg1 == 1 {
				errMsg := "log not defined"
				Service.SetResult(j.ID, j.ExpressionID, nil, &errMsg)
			} else if j.Arg2 <= 0.0 {
				errMsg := "log out of domain"
				Service.SetResult(j.ID, j.ExpressionID, nil, &errMsg)
			} else {
				res := math.Log(j.Arg2) / math.Log(j.Arg1)
				Service.SetResult(j.ID, j.ExpressionID, &res, nil)
			}
		case models.SquareRoot:
			if j.Arg1 < 0 {
				errMsg := "negative square"
				Service.SetResult(j.ID, j.ExpressionID, nil, &errMsg)
			} else {
				res := math.Sqrt(j.Arg1)
				Service.SetResult(j.ID, j.ExpressionID, &res, nil)
			}
		}
	}
}
func TestGetExpression(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestGetExpression struct {
		ID           int64
		Expected     models.Expression
		name         string
		expectedCode int
	}
	tests := []TestGetExpression{
		TestGetExpression{
			name:         "Succes get expression ",
			ID:           1111,
			Expected:     Expr,
			expectedCode: http.StatusOK,
		},
		TestGetExpression{
			name:         "Succes get expression 2",
			ID:           1112,
			Expected:     Expr2,
			expectedCode: http.StatusOK,
		},
		TestGetExpression{
			name:         "bad reqquest",
			ID:           -1,
			expectedCode: http.StatusBadRequest,
		},
		TestGetExpression{
			name:         "not found",
			ID:           math.MaxInt64,
			expectedCode: http.StatusNotFound,
		},
	}
	db.Create(&Expr)
	db.Create(&Expr2)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions/"+fmt.Sprintf("%d", test.ID), nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var expression models.Expression
				data, err := io.ReadAll(w.Result().Body)
				assert.Nil(t, err)
				err = json.Unmarshal(data, &expression)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.ID, expression.ID)
				assert.Equal(t, test.Expected.Status, expression.Status)
				assert.Equal(t, *test.Expected.Result, *expression.Result)
			}
		})
	}
}
func TestGetExpressions(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestGetExpression struct {
		Expected     []models.Expression
		name         string
		expectedCode int
	}
	tests := []TestGetExpression{
		TestGetExpression{
			name: "Succes get expressions ",
			Expected: []models.Expression{
				Expr2,
				Expr,
			},
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&Expr)
	db.Create(&Expr2)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var expressions []models.Expression
				data, err := io.ReadAll(w.Result().Body)
				assert.Nil(t, err)
				err = json.Unmarshal(data, &expressions)
				assert.Nil(t, err)
				for i := 0; i < len(test.Expected); i++ {

					assert.Equal(t, test.Expected[i].Status, expressions[i].Status)
					assert.Equal(t, *test.Expected[i].Result, *expressions[i].Result)
				}
			}
		})
	}
}
func TestGetTask(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestGetTask struct {
		Expected     *models.Task
		name         string
		expectedCode int
	}
	tests := []TestGetTask{
		TestGetTask{
			name: "Succes get task ",
			Expected: &models.Task{
				Arg1:          2.0,
				Arg2:          2.0,
				Operation:     models.Addition,
				OperationTime: config.TimeAddiction.Milliseconds(),
			},
			expectedCode: http.StatusOK,
		},
		TestGetTask{
			name:         "Task not found ",
			expectedCode: http.StatusNotFound,
		},
	}
	for i := 0; i < 10; i++ {
		solveExpr()
		time.Sleep(100 * time.Millisecond)
	}
	Service.CreateExpression("2+2")
	time.Sleep(100 * time.Millisecond)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/internal/task", nil)
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var task *models.Task
				data, err := io.ReadAll(w.Result().Body)
				assert.Nil(t, err)
				err = json.Unmarshal(data, &task)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.Arg1, task.Arg1)
				assert.Equal(t, test.Expected.Arg2, task.Arg2)
				assert.Equal(t, test.Expected.Operation, task.Operation)
				assert.Equal(t, test.Expected.OperationTime, task.OperationTime)
			}
		})
	}
}

type PostTask struct {
	ID           int64   `json:"id"`
	ExpressionID int64   `json:"expression_id"`
	Result       float64 `json:"result"`
}

func TestPostTask(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestPostTask struct {
		ToPost       *PostTask
		name         string
		expectedCode int
	}
	tests := []TestPostTask{
		TestPostTask{
			name:         "Succes post task ",
			ToPost:       &PostTask{ID: 1, ExpressionID: 1, Result: 4.0},
			expectedCode: http.StatusOK,
		},
	}
	for i := 0; i < 10; i++ {
		solveExpr()
		time.Sleep(100 * time.Millisecond)
	}
	Service.CreateExpression("2+2")
	time.Sleep(100 * time.Millisecond)
	task, _ := Service.GetTask()
	tests[0].ToPost.ID = task.ID
	tests[0].ToPost.ExpressionID = task.ExpressionID

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := json.Marshal(test.ToPost)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/internal/task", bytes.NewBuffer(jsonData))
			assert.Nil(t, err)

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			time.Sleep(100 * time.Millisecond)
			var expression models.Expression
			result := db.Model(&models.Expression{}).Where("id = ?", test.ToPost.ExpressionID).First(&expression)
			assert.Nil(t, result.Error)
			assert.Equal(t, test.ToPost.Result, *expression.Result)
		})
	}
}
