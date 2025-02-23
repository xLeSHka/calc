package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
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
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := srv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("server shutdown", zap.Error(err))
				}
			}()
			log.Info("server started at address", zap.String("address", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Error("server shutdown", zap.Error(err))
				return err
			}
			log.Error("server stopped")
			return nil
		},
	})
	return r
}
