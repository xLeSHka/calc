package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) CreateExpression(c *gin.Context) {
	userID, err := jwt.Parse(c)
	if err != nil {
		err.(*customError.CustomError).SendError(c)
		c.Abort()
		return
	}
	var req CreateExpressionReq
	if err := r.validator.ShouldBindJSON(c, &req); err != nil {
		r.Log.Error("CreateExpression: Failed bind body", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("CreateExpression err: %w", err)).Error()
		c.Abort()
		return
	}

	id, cErr := r.usersService.CreateExpression(*req.Expression, userID)
	if cErr != nil {
		r.Log.Error("CreateExpression: Failed create expression", zap.Error(cErr.Err))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated,
		CreateExpressionResp{ID: id},
	)
}
