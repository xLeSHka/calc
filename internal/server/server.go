package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/xLeSHka/calc/pkg/calculator"
	"github.com/xLeSHka/calc/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	server *echo.Echo
}

// создаем новый grpc сервер с gateway`ем
func New(ctx context.Context, port int) (*Server, error) {
	e := echo.New()
	//добавляем мидлвар логирования, recover, CORS
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format:           `{"time":"${time_rfc3339}, "host":"${host}", "method":"${method}", "uri":"${uri}", "status":${status}, "error":"${error}}` + "\n",
			CustomTimeFormat: "2006-01-02 15:04:05",
		},
	))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{"POST", "OPTION"},
	}))
	//создаем экземпляр сервера и добавляем обработчик запросов
	server := &Server{e}
	server.server.POST("/api/v1/calculate", server.calculate)
	return server, nil
}
func (s *Server) calculate(c echo.Context) error {
	//проверка Content-Type в запросе
	if c.Request().Header.Get("Content-Type") != "application/json" {
		logger.New().Info(context.Background(), "content type not allowed", zap.String("Content-Type", c.Request().Header.Get("Content-Type")), zap.String("required Content-Type", "application/json"))
		return echo.NewHTTPError(422, "Expression is not valid")
	}
	//читаем тело запроса и пытаемся его десериализировать
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(422, "Expression is not valid")
	}
	req := Request{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		logger.New().Info(context.Background(), "req", zap.String("expression", err.Error()))

		return echo.NewHTTPError(422, "Expression is not valid")
	}
	//заглушка для 500 ответа
	if req.Expression == "internal" {
		return echo.NewHTTPError(500, "Internal server error")
	}
	//отправляем выражение на вычисление
	res, err := calculator.Calc(req.Expression)
	if err != nil {
		return echo.NewHTTPError(422, "Expression is not valid")
	}
	return c.JSON(http.StatusOK, Response{Result: res, Expression: req.Expression})
}

type Request struct {
	Expression string `json:"expression"`
}
type Response struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

// Функция запуска сервера
func (s *Server) Start(port string) error {
	return s.server.Start(port)
}

// Функция плавного завершения
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
