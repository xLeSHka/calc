package handlers

import (
	"github.com/xLeSHka/calc/internal/orchestrator/transport/routers"
	"go.uber.org/zap"
)

type Router struct {
	Router  *routers.Routers
	service routers.Service
	Log     *zap.Logger
}

func SetUpRouter(
	routers *routers.Routers,
	logger *zap.Logger,
	service routers.Service,
) *Router {
	router := &Router{
		Router:  routers,
		Log:     logger,
		service: service,
	}
	routers.Public.POST("/calculate", router.CreateExpression)
	routers.Public.GET("/expressions", router.GetExpressions)
	routers.Public.GET("/expressions/:id", router.GetExpression)
	routers.Public.GET("/internal/task", router.GetTask)
	routers.Public.POST("/internal/task", router.PostResult)
	return router
}
