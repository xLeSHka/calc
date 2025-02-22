package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET,PUT,DELETE,PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func New(config config.Config, lc fx.Lifecycle, log *zap.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go r.Run(fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort))
			log.Info("server started at address", zap.String("address", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)))
			return nil
		},
		OnStop: nil,
	})
	return r
}
