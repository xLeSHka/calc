package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"net/http"
)

func Parse(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get("userID")
	if !ok {
		return uuid.UUID{}, customError.New(http.StatusUnauthorized, fmt.Errorf("Bad uuid"))
	}
	uid, err := uuid.Parse(id.(string))
	if err != nil {
		return uuid.UUID{}, customError.New(http.StatusUnauthorized, fmt.Errorf("Bad uuid"))
	}
	return uid, nil
}
