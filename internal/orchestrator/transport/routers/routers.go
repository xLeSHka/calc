package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/calc/internal/orchestrator/transport/middlewares"
	"github.com/xLeSHka/calc/internal/utils/jwt"
)

type Routers struct {
	Public      *gin.RouterGroup
	ClientRoute *gin.RouterGroup
}

func CreateRouter(g *gin.Engine, jwt *jwt.JWT, rdb *redis.Client) *Routers {
	public := g.Group("/api/v1")
	private := public.Group("")
	private.Use(middlewares.Auth(jwt, rdb))
	return &Routers{
		Public:      public,
		ClientRoute: private,
	}
}
