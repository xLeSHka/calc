package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) GetExpression(c *gin.Context) {
	userID, err := jwt.Parse(c)
	if err != nil {
		err.(*customError.CustomError).SendError(c)
		c.Abort()
		return
	}
	var req GetExpressionReq
	if err := r.validator.ShouldBindUri(c, &req); err != nil {
		r.Log.Error("GetExpression: Failed bind uri", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("GetExpression: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	expr, cErr := r.usersService.GetExpression(*req.ID, userID)
	if cErr != nil {
		r.Log.Error("GetExpression: Failed get expression", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, mapExpression(expr))
}
