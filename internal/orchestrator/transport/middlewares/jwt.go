package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"net/http"
	"strings"
)

func Auth(jwt *jwt.JWT, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			customError.New(http.StatusUnauthorized, fmt.Errorf("Header not found")).SendError(c)
			c.Abort()
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			customError.New(http.StatusUnauthorized, fmt.Errorf("Header bad format")).SendError(c)
			c.Abort()
			return
		}
		data, err := jwt.VerifyToken(splitToken[1])
		if err != nil {
			customError.New(http.StatusUnauthorized, err).SendError(c)
			c.Abort()
			return
		}

		personID := data["id"].(string)
		val, err := rdb.Get(context.Background(), "jwt:"+personID).Result()
		if err != nil {
			customError.New(http.StatusUnauthorized, err).SendError(c)
			c.Abort()
			return
		} else {
			if val != splitToken[1] {
				customError.New(http.StatusUnauthorized, fmt.Errorf("Token ")).SendError(c)
				c.Abort()
				return
			}
		}
		c.Set("userID", personID)
		c.Next()
	}
}
